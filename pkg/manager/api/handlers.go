package api

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/IPampurin/ImageProcessor/pkg/domain"
	"github.com/IPampurin/ImageProcessor/pkg/manager/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// LoadImageFromProcess возвращает файл изображения по ID и варианту
func LoadImageFromProcess(svc service.ServiceMethods, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Получаем ID из пути
		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
			return
		}

		// 2. Получаем вариант из query (по умолчанию "original")
		variant := c.DefaultQuery("variant", "original")

		// 3. Вызываем сервис
		reader, contentType, err := svc.GetImage(c.Request.Context(), id, variant, log)
		if err != nil {
			log.Error("ошибка получения изображения", "error", err, "id", id, "variant", variant)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer reader.Close() // закрываем после отправки

		// 4. Отдаём файл
		c.DataFromReader(http.StatusOK, -1, contentType, reader, nil)
	}
}

// DeleteImage удаляет изображение и все его варианты.
func DeleteImage(svc service.ServiceMethods, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
			return
		}

		if err := svc.DeleteImage(c.Request.Context(), id, log); err != nil {
			log.Error("ошибка удаления изображения", "error", err, "id", id)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}

// вспомогательные структуры для ответа на /images
type imageVariantResponse struct {
	Type   string `json:"type"`
	Width  *int   `json:"width,omitempty"`
	Height *int   `json:"height,omitempty"`
	Size   int64  `json:"size"`
}

type imageResponse struct {
	ID           uuid.UUID              `json:"id"`
	OriginalName string                 `json:"originalName"`
	Status       string                 `json:"status"`
	Variants     []imageVariantResponse `json:"variants"`
}

// GetImages возвращает список последних загруженных изображений.
func GetImages(svc service.ServiceMethods, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Получаем limit из query (по умолчанию 20)
		limitStr := c.DefaultQuery("limit", "20")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 20
		}

		// 2. Вызываем сервис (получаем доменные модели)
		images, err := svc.ListImages(c.Request.Context(), limit, log)
		if err != nil {
			log.Error("ошибка получения списка изображений", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 3. Преобразуем в ответ для фронта
		response := make([]imageResponse, 0, len(images))
		for _, img := range images {
			// Здесь нужно собрать варианты. Пока variants пустой, т.к. ещё не реализовано.
			// В будущем можно будет заполнять из дочерних записей.
			resp := imageResponse{
				ID:           img.ID,
				OriginalName: img.Name,
				Status:       img.Status,
				Variants:     []imageVariantResponse{}, // пока пусто
			}
			response = append(response, resp)
		}

		c.JSON(http.StatusOK, response)
	}
}
