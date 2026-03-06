package db

import (
	"time"

	"github.com/google/uuid"
)

// ImageFile представляет запись в таблице image_files,
// хранящей как исходные изображения, так и все их обработанные варианты.
type ImageFile struct {
	ID           uuid.UUID  `db:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"` // уникальный идентификатор файла
	OriginalID   *uuid.UUID `db:"original_id" gorm:"index"`                                 // ссылка на исходный файл (NULL для исходного)
	Name         string     `db:"name"`                                                     // полное имя файла (с суффиксом)
	Type         string     `db:"type"`                                                     // тип: 'original', 'resized', 'thumbnail', 'watermarked'
	ContentType  string     `db:"content_type"`                                             // MIME-тип (image/jpeg, image/png)
	Size         int64      `db:"size"`                                                     // размер файла в байтах
	Width        *int       `db:"width"`                                                    // ширина (если есть)
	Height       *int       `db:"height"`                                                   // высота (если есть)
	Status       string     `db:"status"`                                                   // статус обработки (для исходного: pending/processing/completed/failed; для вариантов: completed)
	ErrorMessage *string    `db:"error_message"`                                            // текст ошибки (только для исходного при статусе failed)
	StoragePath  string     `db:"storage_path"`                                             // ключ в S3
	CreatedAt    time.Time  `db:"created_at"`                                               // время попадания в S3
}

// Outbox представляет сообщение для гарантированной отправки в Kafka (Transactional Outbox pattern).
type Outbox struct {
	ID        uuid.UUID `db:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"` // уникальный идентификатор сообщения
	Topic     string    `db:"topic"`                                                    // топик Kafka ('image-tasks')
	Key       string    `db:"key"`                                                      // ключ сообщения (обычно image_id)
	Payload   []byte    `db:"payload"`                                                  // данные задачи (сериализованный Task)
	CreatedAt time.Time `db:"created_at"`                                               // время создания записи
}
