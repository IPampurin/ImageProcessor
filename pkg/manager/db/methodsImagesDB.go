package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// InsertImage создаёт запись в images
func (d *DataBase) InsertImage(ctx context.Context, i *Image) error {

	query := `INSERT INTO images (id, original_id, name, type, content_type, size, width, height, status, error_message, storage_path, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := d.Pool.Exec(ctx, query,
		i.ID,
		i.OriginalID,
		i.Name,
		i.Type,
		i.ContentType,
		i.Size,
		i.Width,
		i.Height,
		i.Status,
		i.ErrorMessage,
		i.StoragePath,
		i.CreatedAt)
	if err != nil {
		return fmt.Errorf("ошибка InsertImage добавления записи в images: %w", err)
	}

	return nil
}

// GetByID возвращает запись из images по уникальному идентификатору
func (d *DataBase) GetByID(ctx context.Context, uid uuid.UUID) (*Image, error) {

	query := `SELECT *
	            FROM images
			   WHERE id = $1`

	i := &Image{}

	err := d.Pool.QueryRow(ctx, query, uid).Scan(
		&i.ID,
		&i.OriginalID,
		&i.Name,
		&i.Type,
		&i.ContentType,
		&i.Size,
		&i.Width,
		&i.Height,
		&i.Status,
		&i.ErrorMessage,
		&i.StoragePath,
		&i.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("ошибка GetByID получения записи в images: %w", err)
	}

	return i, nil
}

// UpdateStatusOrErr обновляет запсь в images по статусу или ошибке
func (d *DataBase) UpdateStatusOrErr(ctx context.Context, uid uuid.UUID, status string, errMsg *string) error {

	query := `UPDATE images
	             SET status = $1, error_message = $2
			   WHERE id = $3`

	_, err := d.Pool.Exec(ctx, query, status, errMsg, uid)
	if err != nil {
		return fmt.Errorf("ошибка UpdateStatusOrErr при изменении записи в images: %w", err)
	}

	return nil
}

// DeleteImage удаляет запись из images
func (d *DataBase) DeleteImage(ctx context.Context, uid uuid.UUID) error {

	query := `DELETE
	            FROM images
	           WHERE id = $1`

	_, err := d.Pool.Exec(ctx, query, uid)
	if err != nil {
		return fmt.Errorf("ошибка DeleteImage при удалении записи из images: %w", err)
	}

	return nil
}

// ListLatestOriginals используется для отображения UI изображений в галерее
func (d *DataBase) ListLatestOriginals(ctx context.Context, limit int) ([]*Image, error) {

	query := `SELECT *
	            FROM images
			   WHERE original_id IS NULL
			   ORDER BY created_at DESC
			   LIMIT $1`

	rows, err := d.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка ListLatestOriginals при получении записей из images: %w", err)
	}
	defer rows.Close()

	result := make([]*Image, 0)
	for rows.Next() {
		var i Image
		err := rows.Scan(
			&i.ID,
			&i.OriginalID,
			&i.Name,
			&i.Type,
			&i.ContentType,
			&i.Size,
			&i.Width,
			&i.Height,
			&i.Status,
			&i.ErrorMessage,
			&i.StoragePath,
			&i.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка ListLatestOriginals при сканировании записи из images: %w", err)
		}

		result = append(result, &i)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка ListLatestOriginals при итерации по записям из images: %w", err)
	}

	return result, nil
}
