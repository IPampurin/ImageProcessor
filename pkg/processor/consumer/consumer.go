package consumer

import (
	"context"
	"strings"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/logger"
)

// Consumer обёртка над kafka.Consumer из фреймворка wbf
type Consumer struct {
	reader *kafka.Consumer
	log    logger.Logger
}

// New создаёт нового потребителя Kafka
func New(brokers string, topic, groupID string, log logger.Logger) *Consumer {
	return &Consumer{
		reader: kafka.NewConsumer(strings.Split(brokers, ","), topic, groupID),
		log:    log,
	}
}

// StartConsuming запускает фоновое чтение сообщений из Kafka.
// Внимание: для полной аналогии с менеджером сюда нужно добавить параметр стратегии повторов.
func (c *Consumer) StartConsuming(ctx context.Context, tasksCh chan<- kafkago.Message) {
	go func() {
		for {
			msg, err := c.reader.Fetch(ctx)
			if err != nil {
				c.log.Error("ошибка получения сообщения из Kafka", "error", err)
				return
			}
			select {
			case tasksCh <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Commit коммитит сообщение через reader.Commit (именно так, как в менеджере)
func (c *Consumer) Commit(ctx context.Context, msg kafkago.Message) error {
	// ВАЖНО: используем метод Commit, а не CommitMessages
	return c.reader.Commit(ctx, msg)
}

// Close закрывает соединение с Kafka
func (c *Consumer) Close() error {
	return c.reader.Close()
}
