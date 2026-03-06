package s3

import (
	"context"
	"io"
)

// StorageS3 описывает методы для работы с распределённым хранилищем
type StorageS3 interface {

	// Upload сохраняет файл в хранилище по указанному ключу (пути)
	Upload(ctx context.Context, path string, reader io.Reader, size int64, contentType string) error

	// Download возвращает ReadCloser для чтения файла по указанному ключу (пути)
	Download(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete удаляет файл по указанному ключу (по пути)
	Delete(ctx context.Context, path string) error
}
