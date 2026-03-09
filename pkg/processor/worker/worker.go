package worker

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"path/filepath"

	"github.com/IPampurin/ImageProcessor/pkg/processor/models"
	"github.com/IPampurin/ImageProcessor/pkg/processor/operations"
	"github.com/IPampurin/ImageProcessor/pkg/s3"
	"github.com/wb-go/wbf/logger"
)

// Worker выполняет операции над изображениями
type Worker struct {
	s3          *s3.S3
	log         logger.Logger
	thumbWidth  int
	thumbHeight int
}

// New создаёт нового воркера
func New(s3 *s3.S3, log logger.Logger, thumbWidth, thumbHeight int) *Worker {
	return &Worker{
		s3:          s3,
		log:         log,
		thumbWidth:  thumbWidth,
		thumbHeight: thumbHeight,
	}
}

// ProcessTask обрабатывает одну задачу: скачивает оригинал, применяет операции, загружает результаты
func (w *Worker) ProcessTask(ctx context.Context, task *models.Task) (*models.Result, error) {
	// Скачиваем оригинал из S3
	reader, err := w.s3.Download(ctx, task.ObjectKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка скачивания оригинала из S3: %w", err)
	}
	defer reader.Close()

	// Декодируем изображение
	img, format, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования изображения: %w", err)
	}

	// Определяем MIME-тип
	contentType := "image/jpeg"
	if format == "png" {
		contentType = "image/png"
	}

	var variants []models.VariantResult

	// Операция: ресайз
	if task.Resize != nil {
		resizedImg, err := operations.Resize(img, task.Resize.Width, task.Resize.Height)
		if err != nil {
			w.log.Error("ошибка ресайза", "error", err, "imageID", task.ImageID)
		} else {
			key := generateKey(task.ImageID, "resized", filepath.Ext(task.ObjectKey))
			if err := w.uploadImage(ctx, key, resizedImg, contentType); err != nil {
				w.log.Error("ошибка загрузки ресайза в S3", "error", err, "key", key)
			} else {
				width := task.Resize.Width
				height := task.Resize.Height
				variants = append(variants, models.VariantResult{
					Type:        "resized",
					StoragePath: key,
					ContentType: contentType,
					Size:        0, // TODO: получить фактический размер после загрузки
					Width:       &width,
					Height:      &height,
				})
			}
		}
	}

	// Операция: миниатюра
	if task.Thumbnail {
		thumbImg, err := operations.Thumbnail(img, w.thumbWidth, w.thumbHeight)
		if err != nil {
			w.log.Error("ошибка создания миниатюры", "error", err, "imageID", task.ImageID)
		} else {
			key := generateKey(task.ImageID, "thumbnail", filepath.Ext(task.ObjectKey))
			if err := w.uploadImage(ctx, key, thumbImg, contentType); err != nil {
				w.log.Error("ошибка загрузки миниатюры в S3", "error", err, "key", key)
			} else {
				width := w.thumbWidth
				height := w.thumbHeight
				variants = append(variants, models.VariantResult{
					Type:        "thumbnail",
					StoragePath: key,
					ContentType: contentType,
					Size:        0,
					Width:       &width,
					Height:      &height,
				})
			}
		}
	}

	// Операция: водяной знак (пока заглушка, чтобы не требовала шрифт)
	if task.Watermark {
		watermarkedImg, err := operations.AddTextWatermark(img, "Sample Watermark")
		if err != nil {
			w.log.Error("ошибка наложения водяного знака", "error", err, "imageID", task.ImageID)
		} else {
			key := generateKey(task.ImageID, "watermarked", filepath.Ext(task.ObjectKey))
			if err := w.uploadImage(ctx, key, watermarkedImg, contentType); err != nil {
				w.log.Error("ошибка загрузки изображения с водяным знаком в S3", "error", err, "key", key)
			} else {
				bounds := watermarkedImg.Bounds()
				width := bounds.Dx()
				height := bounds.Dy()
				variants = append(variants, models.VariantResult{
					Type:        "watermarked",
					StoragePath: key,
					ContentType: contentType,
					Size:        0,
					Width:       &width,
					Height:      &height,
				})
			}
		}
	}

	// Формируем результат
	result := &models.Result{
		ImageID:  task.ImageID,
		Status:   "completed",
		Variants: variants,
	}

	// Если ни одна операция не удалась, помечаем как failed
	if len(variants) == 0 && (task.Resize != nil || task.Thumbnail || task.Watermark) {
		errMsg := "все операции завершились ошибкой"
		result.Status = "failed"
		result.ErrorMessage = &errMsg
	}

	return result, nil
}

// uploadImage кодирует изображение в JPEG/PNG и загружает в S3
func (w *Worker) uploadImage(ctx context.Context, key string, img image.Image, contentType string) error {
	var data []byte
	var err error

	if contentType == "image/jpeg" {
		data, err = operations.EncodeJPEG(img)
	} else {
		data, err = operations.EncodePNG(img)
	}
	if err != nil {
		return fmt.Errorf("ошибка кодирования изображения: %w", err)
	}

	return w.s3.Upload(ctx, key, bytes.NewReader(data), contentType)
}

// generateKey формирует ключ для S3 вида: {imageID}/{variant}.ext
func generateKey(imageID, variant, ext string) string {
	return fmt.Sprintf("%s/%s%s", imageID, variant, ext)
}
