package controllers

import (
	"EDU_TH_2_backend/gin/logger"
	"EDU_TH_2_backend/gin/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"image/png"
	"net/http"
)

type UtilController struct {
	utilService services.UtilService
}

func NewUtilController(utilService services.UtilService) *UtilController {
	return &UtilController{utilService: utilService}
}

func (ctrl *UtilController) GetQRCode(c *gin.Context) {
	// 获取参数值
	url := c.Query("url")

	qrCode, err := ctrl.utilService.GetQRCode(url)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "生成二维码失败")
	}

	// 设置响应头，告诉客户端返回的是图片
	c.Header("Content-Type", "image/png")

	// 输出为 PNG 图片
	png.Encode(c.Writer, qrCode)
}

func (ctrl *UtilController) ImgSave(c *gin.Context) {

	// 获取body的from-data中的file类型文件
	file, err := c.FormFile("data")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "获取文件失败")
		return
	}
	//// 保存文件到本地
	url, err := ctrl.utilService.ImgSave(c, file)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("保存文件失败：%s", err.Error()))
		return
	}

	// log:哪个用户干了些什么（从auth中获取用户信息）
	logger.Info("用户上传了图片")

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}
