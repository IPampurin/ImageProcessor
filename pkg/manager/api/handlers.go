package api

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/IPampurin/ImageProcessor/pkg/domain"
	"github.com/IPampurin/ImageProcessor/pkg/manager/service"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/logger"
)

// UploadImageToProcess возвращает gin.HandlerFunc для загрузки изображения
func UploadImageToProcess(svc service.ServiceMethods, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Получаем файл из формы
		file, header, err := c.Request.FormFile("image")
		if err != nil {
			log.Error("не удалось получить файл", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "файл не найден"})
			return
		}
		defer file.Close() // закроем после того, как сервис прочитает данные

		// 2. Определяем MIME-тип по первым 512 байтам
		buff := make([]byte, 512)
		n, err := file.Read(buff)
		if err != nil && err != io.EOF {
			log.Error("ошибка чтения первых байт", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось прочитать файл"})
			return
		}
		contentType := http.DetectContentType(buff[:n])

		// 3. Создаём MultiReader, который вернёт прочитанные байты, а затем остаток файла
		//    (это позволяет не терять данные, которые мы уже прочитали для определения MIME)
		reader := io.MultiReader(bytes.NewReader(buff[:n]), file)

		// 4. Парсим остальные поля формы
		thumbnail := c.PostForm("thumbnail") == "true"
		watermark := c.PostForm("watermark") == "true"

		var resize *domain.ResizeOptions
		widthStr := c.PostForm("width")
		heightStr := c.PostForm("height")

		if widthStr != "" || heightStr != "" {
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
			resize = &domain.ResizeOptions{Width: width, Height: height}
		}

		// 5. Формируем доменную структуру запроса (без HTTP-зависимостей)
		uploadData := &domain.UploadData{
			Filename:    header.Filename,
			ContentType: contentType,
			Size:        header.Size,
			Reader:      reader,
			Thumbnail:   thumbnail,
			Watermark:   watermark,
			Resize:      resize,
		}

		// 6. Вызываем сервис
		id, err := svc.UploadImage(c.Request.Context(), uploadData, log)
		if err != nil {
			log.Error("ошибка загрузки изображения", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 7. Возвращаем UUID
		c.JSON(http.StatusAccepted, UploadResponse{ID: id})
	}
}
