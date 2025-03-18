package utils

import (
	"EDU_TH_2_backend/gin/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"time"
)

var minioClient *minio.Client

func InitMinio() error {
	endpoint := config.GetString("minio.endpoint")
	accessKeyID := config.GetString("minio.access_key_id")
	secretAccessKey := config.GetString("minio.secret_access_key")
	useSSL := false

	// 初始化 MinIO 客户端
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}

	minioClient = client

	return nil
}

func GetMinioClient() *minio.Client {
	return minioClient
}

func GetAvatarURL(c *gin.Context, userID string) (string, error) {
	filename := "avatar.png"

	//avatarPath := filepath.Join(".", "gin", "uploads", "avatars", userID+"-"+filename)

	//_, err := os.Stat(avatarPath)
	//if os.IsNotExist(err) {
	//	// 如果用户没有头像，则返回默认头像
	//	defaultAvatarPath := filepath.Join(".", "gin", "static", "default_avatar.png")
	//	return defaultAvatarPath, nil
	//} else {
	//	// 如果用户有头像，则返回该头像
	//	return avatarPath, nil
	//}

	bucketName := fmt.Sprintf("user-%s", userID)

	// 桶中是否存在该文件
	_, err := minioClient.StatObject(c, bucketName, filename, minio.StatObjectOptions{})
	if err != nil {
		// 如果用户没有头像，则返回默认头像（default-bucket桶中的default_avatar.png）
		object, err := minioClient.PresignedGetObject(c, "default-bucket", "default_avatar.png", 7*24*time.Hour, nil)
		if err != nil {
			return "", err
		}

		return object.String(), nil
	} else {
		// 如果用户有头像，则返回该头像
		object, err := minioClient.PresignedGetObject(c, bucketName, filename, 7*24*time.Hour, nil)
		if err != nil {
			return "", err
		}

		return object.String(), nil
	}
}
