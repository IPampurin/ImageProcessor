package db

import (
	"context"

	"github.com/IPampurin/ImageProcessor/pkg/manager/db"
	"github.com/google/uuid"
)

type StorageDB interface {
	ImageFileMethods
	OutboxMethods
}

// ImageFileMethods - интерфейс с методами для работы с таблицей image_files внутреннего хранилища
type ImageFileMethods interface {

	// CreateOriginal создаёт запись в image_files при загрузке изображения в обрабоку (статус pending)
	CreateOriginal(ctx context.Context, name, contentType string, size int64) (*db.ImageFile, error)

	// GetByID возвращает запись из image_files по уникальному идентификатору
	GetByID(ctx context.Context, id uuid.UUID) (*db.ImageFile, error)

	// UpdateStatusOrErr обновляет запсь в image_files по статусу или ошибке
	UpdateStatusOrErr(ctx context.Context, id uuid.UUID, status string, errMsg *string) error

	// DeleteImage удаляет запись из image_files
	DeleteImage(ctx context.Context, id uuid.UUID) error
}

// OutboxMethods - интерфейс с методами для работы с таблицей outbox внутреннего хранилища
type OutboxMethods interface {

	// CreateOutbox создаёт запись в таблице outbox
	CreateOutbox(ctx context.Context, topic, key string, payload []byte) error

	// GetUnsentOutbox получает свежие записи для отправки брокеру
	GetUnsentOutbox(ctx context.Context, limit int) ([]db.Outbox, error)

	// DeleteOutbox удаляет запись из таблицы outbox
	DeleteOutbox(ctx context.Context, id uuid.UUID) error
}
