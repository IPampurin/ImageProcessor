package service

import (
	"context"
	"io"

	"github.com/IPampurin/ImageProcessor/pkg/domain"
	"github.com/IPampurin/ImageProcessor/pkg/manager/db"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/logger"
)

// ServiceMethods
type ServiceMethods interface {

	// UploadImage загружает изображение, сохраняет его в S3, создаёт запись в БД и задачу в outbox
	UploadImage(ctx context.Context, data *domain.UploadData, log logger.Logger) (uuid.UUID, error)

	// GetImage возвращает файл изображения по его ID и варианту (original, thumbnail, resized и т.д.).
	// Возвращает ReadCloser (который нужно закрыть после использования), ContentType и ошибку.
	GetImage(ctx context.Context, id uuid.UUID, variant string, log logger.Logger) (io.ReadCloser, string, error)

	// DeleteImage удаляет изображение и все его обработанные варианты из БД и S3.
	DeleteImage(ctx context.Context, id uuid.UUID, log logger.Logger) error

	// ListImages возвращает список последних загруженных оригинальных изображений
	// для отображения в галерее (с пагинацией через limit).
	ListImages(ctx context.Context, limit int, log logger.Logger) ([]*db.Image, error)
}
