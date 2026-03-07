package api

import (
	"net/http"
	"strconv"

	"github.com/IPampurin/ImageProcessor/pkg/manager/service"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/logger"
)

// UploadImageToProcess возвращает gin.HandlerFunc для загрузки изображения
func UploadImageToProcess(svc service.ServiceMethods, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		// получаем файл из формы
		file, header, err := c.Request.FormFile("image")
		if err != nil {
			log.Error("не удалось получить файл", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "файл не найден"})
			return
		}
		defer file.Close()

		// парсим остальные поля
		var req UploadRequest
		req.File = file
		req.FileHeader = header
		req.Thumbnail = c.PostForm("thumbnail") == "true"
		req.Watermark = c.PostForm("watermark") == "true"

		// обрабатываем параметры ресайза
		widthStr := c.PostForm("width")
		heightStr := c.PostForm("height")
		if widthStr != "" || heightStr != "" {
			// если одно из полей указано, должны быть оба
			if widthStr == "" || heightStr == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "необходимо указать и ширину, и высоту для ресайза"})
				return
			}
			width, err := strconv.Atoi(widthStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ширина должна быть числом"})
				return
			}
			height, err := strconv.Atoi(heightStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "высота должна быть числом"})
				return
			}
			req.Resize = &ResizeOptions{Width: width, Height: height}
		}

		// вызываем сервис
		id, err := svc.UploadImage(c.Request.Context(), &req)
		if err != nil {
			log.Error("ошибка загрузки изображения", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// возвращаем UUID
		c.JSON(http.StatusAccepted, UploadResponse{ID: id})
	}
}
