package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IPampurin/ImageProcessor/pkg/configuration"
	"github.com/IPampurin/ImageProcessor/pkg/processor/consumer"
	"github.com/IPampurin/ImageProcessor/pkg/processor/models"
	"github.com/IPampurin/ImageProcessor/pkg/processor/producer"
	"github.com/IPampurin/ImageProcessor/pkg/processor/worker"
	"github.com/IPampurin/ImageProcessor/pkg/s3"
	"github.com/wb-go/wbf/logger"

	kafkago "github.com/segmentio/kafka-go"
)

func main() {
	// Создаём корневой контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем горутину обработки сигналов ОС
	go signalHandler(ctx, cancel)

	// Загружаем конфигурацию из .env
	cfg, err := configuration.ReadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализируем логгер (точно как в менеджере)
	appLogger, err := logger.InitLogger(
		logger.ZapEngine,
		"ImageProcessor",
		os.Getenv("APP_ENV"),
		logger.WithLevel(logger.InfoLevel),
	)
	if err != nil {
		log.Fatalf("Ошибка создания логгера: %v", err)
	}
	// defer в точности как в менеджере: без Sync, просто приведение типа
	defer func() { _ = appLogger.(*logger.ZapAdapter) }()

	// Подключаемся к S3 (MinIO)
	s3Client, err := s3.InitS3(ctx, &cfg.S3, appLogger)
	if err != nil {
		appLogger.Error("ошибка подключения к S3 хранилищу", "error", err)
		return
	}

	// Создаём Kafka consumer и producer
	kafkaConsumer := consumer.New(cfg.Kafka.Brokers, cfg.Kafka.InputTopic, cfg.Kafka.ConsumerGroup, appLogger)
	defer kafkaConsumer.Close()

	kafkaProducer := producer.New(cfg.Kafka.Brokers, cfg.Kafka.OutputTopic, appLogger)
	defer kafkaProducer.Close()

	// Каналы для пайплайна
	tasksCh := make(chan kafkago.Message)
	resultsCh := make(chan *processorResult)

	// Запускаем consumer в фоне
	kafkaConsumer.StartConsuming(ctx, tasksCh)

	// Создаём воркер
	imageWorker := worker.New(s3Client, appLogger, cfg.Thumb.ThumbWidth, cfg.Thumb.ThumbHeight)

	// Запускаем пул воркеров
	const numWorkers = 5
	for i := 0; i < numWorkers; i++ {
		go startWorker(ctx, imageWorker, tasksCh, resultsCh, appLogger)
	}

	// Запускаем отправитель результатов
	go startResultSender(ctx, kafkaConsumer, kafkaProducer, resultsCh, appLogger)

	appLogger.Info("Процессор обработки изображений успешно запущен")

	// Ожидаем сигнала завершения
	<-ctx.Done()
	appLogger.Info("Завершение работы процессора")
	time.Sleep(2 * time.Second)
}

// processorResult связывает исходное сообщение Kafka и результат его обработки
type processorResult struct {
	rawMsg kafkago.Message
	result *models.Result
	err    error
}

// startWorker читает задачи из tasksCh, обрабатывает и отправляет результаты в resultsCh
func startWorker(ctx context.Context, w *worker.Worker, tasksCh <-chan kafkago.Message, resultsCh chan<- *processorResult, log logger.Logger) {
	for {
		select {
		case <-ctx.Done():
			return
		case rawMsg := <-tasksCh:
			var task models.Task
			if err := json.Unmarshal(rawMsg.Value, &task); err != nil {
				log.Error("ошибка разбора задачи JSON", "error", err, "offset", rawMsg.Offset)
				resultsCh <- &processorResult{rawMsg: rawMsg, err: err}
				continue
			}

			log.Info("обработка задачи", "imageID", task.ImageID, "offset", rawMsg.Offset)
			result, err := w.ProcessTask(ctx, &task)
			resultsCh <- &processorResult{rawMsg: rawMsg, result: result, err: err}
		}
	}
}

// startResultSender получает результаты, отправляет их в Kafka и коммитит исходные сообщения
func startResultSender(ctx context.Context, c *consumer.Consumer, p *producer.Producer, resultsCh <-chan *processorResult, log logger.Logger) {

	for {
		select {
		case <-ctx.Done():
			return
		case pres := <-resultsCh:
			if pres.err == nil && pres.result != nil {
				if err := p.SendResult(ctx, pres.result.ImageID, pres.result); err != nil {
					log.Error("ошибка отправки результата в Kafka", "error", err, "imageID", pres.result.ImageID)
					continue
				}
			}
			// Коммитим через consumer (как в менеджере)
			if err := c.Commit(ctx, pres.rawMsg); err != nil {
				log.Error("ошибка коммита сообщения Kafka", "error", err, "offset", pres.rawMsg.Offset)
			}
		}
	}
}

// signalHandler обрабатывает сигналы отмены (как в менеджере)
func signalHandler(ctx context.Context, cancel context.CancelFunc) {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	select {
	case <-ctx.Done():
		return
	case <-sigChan:
		cancel()
		return
	}
}
