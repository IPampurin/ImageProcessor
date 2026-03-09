package producer

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/IPampurin/ImageProcessor/pkg/processor/models"
	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/logger"
)

// Producer обёртка над kafka.Producer из фреймворка wbf
type Producer struct {
	writer *kafka.Producer
	log    logger.Logger
}

// New создаёт нового продюсера Kafka
func New(brokers string, topic string, log logger.Logger) *Producer {
	return &Producer{
		writer: kafka.NewProducer(strings.Split(brokers, ","), topic),
		log:    log,
	}
}

// SendResult отправляет результат обработки в Kafka
func (p *Producer) SendResult(ctx context.Context, key string, result *models.Result) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return p.writer.Send(ctx, []byte(key), data)
}

// Close закрывает соединение с Kafka
func (p *Producer) Close() error {
	return p.writer.Close()
}
