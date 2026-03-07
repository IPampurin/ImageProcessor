package db

import (
	"context"

	"github.com/google/uuid"
)

type StorageDB interface {
	ImageFileMethods
	OutboxMethods
}

// ImageFileMethods - интерфейс с методами для работы с таблицей images внутреннего хранилища
type ImageFileMethods interface {

	// InsertImage создаёт запись в images
	InsertImage(ctx context.Context, i *Image) error

	// GetByID возвращает запись из images по уникальному идентификатору
	GetByID(ctx context.Context, uid uuid.UUID) (*Image, error)

	// UpdateStatusOrErr обновляет запсь в images по статусу или ошибке
	UpdateStatusOrErr(ctx context.Context, id uuid.UUID, status string, errMsg *string) error

	// DeleteImage удаляет запись из images
	DeleteImage(ctx context.Context, uid uuid.UUID) error

	// ListLatestOriginals используется для отображения UI изображений в галерее
	ListLatestOriginals(ctx context.Context, limit int) ([]*Image, error)
}

// OutboxMethods - интерфейс с методами для работы с таблицей outbox внутреннего хранилища
type OutboxMethods interface {

	// CreateOutbox создаёт запись в таблице outbox
	CreateOutbox(ctx context.Context, outRow *Outbox) error

	// GetUnsentOutbox получает свежие записи для отправки брокеру
	GetUnsentOutbox(ctx context.Context, limit int) ([]*Outbox, error)

	// DeleteOutbox удаляет запись из таблицы outbox
	DeleteOutbox(ctx context.Context, uid uuid.UUID) error
}
