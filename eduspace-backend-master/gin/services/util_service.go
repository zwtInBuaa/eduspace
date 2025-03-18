package services

import (
	"EDU_TH_2_backend/gin/utils"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
)

type UtilService interface {
	GetQRCode(url string) (barcode.Barcode, error)
	ImgSave(c *gin.Context, file *multipart.FileHeader) (url string, err error)
}

type utilService struct {
}

func NewUtilService() UtilService {
	return &utilService{}
}

func (s *utilService) GetQRCode(url string) (barcode.Barcode, error) {

	// 生成二维码
	qrCode, err := qr.Encode(url, qr.L, qr.Auto)
	if err != nil {
		return nil, err
	}

	// 转换为图片格式
	qrCode, _ = barcode.Scale(qrCode, 256, 256)

	return qrCode, nil
}

// ImgSave 保存上传的图片
func (s *utilService) ImgSave(c *gin.Context, file *multipart.FileHeader) (url string, err error) {

	// 如果默认桶不存在，则创建默认桶
	minioClient := utils.GetMinioClient()

	bucketName := "default-bucket"

	exists, err := minioClient.BucketExists(c, bucketName)
	if err != nil {
		return "", err
	}

	if !exists {
		err = minioClient.MakeBucket(c, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", err
		}
	}

	// 获取上传的文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 生成一个随机的文件名（注意原有的filename包括了文件后缀）
	filename := uuid.New().String() + "-" + file.Filename

	//uploadInfo, err := minioClient.PutObject(c, "default-bucket", file.Filename, src, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	_, err = minioClient.PutObject(c, "default-bucket", filename, src, -1, minio.PutObjectOptions{ContentType: "image/png"})
	if err != nil {
		return "", err
	}

	// 生成一个有效期为 10 分钟的签名 URL
	//object, err := minioClient.PresignedGetObject(c, bucketName, file.Filename, 7*24*time.Hour, nil)
	//if err != nil {
	//	return "", err
	//}
	objectURL := minioClient.EndpointURL().String() + "/" + bucketName + "/" + filename

	return objectURL, nil

	//return uploadInfo.Location, nil
}

//func (s *utilService) ImgSave(file *multipart.FileHeader) (url string, err error) {
//
//	// 获取上传的文件
//	src, err := file.Open()
//	if err != nil {
//		return "", err
//	}
//	defer src.Close()
//
//	// 创建保存文件的目录
//	dir := "./gin/static"
//	if _, err := os.Stat(dir); os.IsNotExist(err) {
//		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
//			return "", err
//		}
//	}
//
//	// 生成保存文件的路径和文件名
//	filename := filepath.Base(file.Filename)
//	dstPath := filepath.Join(dir, filename)
//
//	// 创建保存文件的目标文件
//	dst, err := os.Create(dstPath)
//	if err != nil {
//		return "", err
//	}
//	defer dst.Close()
//
//	// 将上传的文件内容复制到目标文件中
//	if _, err := io.Copy(dst, src); err != nil {
//		return "", err
//	}
//
//	// 返回保存文件的 URL
//	url = "/static/" + filename
//	return url, nil
//}
