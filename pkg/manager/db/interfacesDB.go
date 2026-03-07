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

// ImageFileMethods - интерфейс с методами для работы с таблицей images внутреннего хранилища
type ImageFileMethods interface {

	// InsertImage создаёт запись в images
	InsertImage(ctx context.Context, image *db.Image) error

	// GetByID возвращает запись из images по уникальному идентификатору
	GetByID(ctx context.Context, id uuid.UUID) (*db.Image, error)

	// UpdateStatusOrErr обновляет запсь в images по статусу или ошибке
	UpdateStatusOrErr(ctx context.Context, id uuid.UUID, status string, errMsg *string) error

	// DeleteImage удаляет запись из images
	DeleteImage(ctx context.Context, id uuid.UUID) error

	// ListLatestOriginals используется для отображения UI изображений в галерее
	ListLatestOriginals(ctx context.Context, limit int) ([]db.Image, error)
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
