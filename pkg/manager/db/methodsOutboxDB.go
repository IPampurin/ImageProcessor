package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CreateOutbox создаёт запись в таблице outbox
func (d *DataBase) CreateOutbox(ctx context.Context, rowToSend *Outbox) error {

	query := `INSERT INTO outbox (id, topic, key, payload, created_at)
	          VALUES ($1, $2, $3, $4, $5)`

	_, err := d.Pool.Exec(ctx, query,
		rowToSend.ID,
		rowToSend.Topic,
		rowToSend.Key,
		rowToSend.Payload,
		rowToSend.CreatedAt)
	if err != nil {
		return fmt.Errorf("ошибка CreateOutbox при добавлении записи в outbox: %w", err)
	}

	return nil
}

// GetUnsentOutbox получает limit записей для отправки брокеру
func (d *DataBase) GetUnsentOutbox(ctx context.Context, limit int) ([]*Outbox, error) {

	query := `SELECT *
	            FROM outbox
			   ORDER BY created_at
			   LIMIT $1`

	rows, err := d.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка GetUnsentOutbox при получении записей из outbox: %w", err)
	}
	defer rows.Close()

	result := make([]*Outbox, 0)
	for rows.Next() {
		var outRow Outbox
		err := rows.Scan(
			&outRow.ID,
			&outRow.Topic,
			&outRow.Key,
			&outRow.Payload,
			&outRow.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка GetUnsentOutbox при сканировании записи из outbox: %w", err)
		}

		result = append(result, &outRow)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка GetUnsentOutbox при итерации по записям из outbox: %w", err)
	}

	return result, nil
}

// DeleteOutbox удаляет запись из таблицы outbox
func (d *DataBase) DeleteOutbox(ctx context.Context, uid uuid.UUID) error {

	query := `DELETE
	            FROM outbox
			   WHERE id = $1`

	_, err := d.Pool.Exec(ctx, query, uid)
	if err != nil {
		return fmt.Errorf("ошибка DeleteOutbox при удалении записи из outbox: %w", err)
	}

	return nil
}
