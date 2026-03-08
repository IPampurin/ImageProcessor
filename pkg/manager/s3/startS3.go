package s3

import (
	"context"
	"fmt"

	"github.com/IPampurin/ImageProcessor/pkg/manager/configuration"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/wb-go/wbf/logger"
)

// S3 хранит подключение к S3-совместимому хранилищу
type S3 struct {
	Client *s3.Client
	Bucket string
}

// InitS3 инициализирует подключение к S3-совместимому хранилищу (например, MinIO)
func InitS3(ctx context.Context, cfg *configuration.ConfS3, log logger.Logger) (*S3, error) {

	// определяем схему (http/https) в соответствии с флагом UseSSL
	scheme := "http"
	if cfg.UseSSL {
		scheme = "https"
	}
	endpoint := fmt.Sprintf("%s://%s", scheme, cfg.Endpoint)

	// загружаем конфигурацию AWS с использованием статических учётных данных
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"), // регион не важен для MinIO
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKey,
			cfg.SecretKey,
			"", // токен сессии не требуется
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка InitS3 определения конфигурации внешнего хранилища: %w", err)
	}

	// создаём клиент S3 с кастомным endpoint'ом и включённым path-style (для MinIO обязательно UsePathStyle = true)
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	// проверяем доступность бакета, указанного в конфигурации
	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(cfg.Bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки бакета: bucket %q не доступен: %w", cfg.Bucket, err)
	}

	log.Info("S3 клиент успешно инициализирован, бакет доступен.")

	return &S3{
		Client: client,
		Bucket: cfg.Bucket}, nil
}

// CloseS3 освобождает ресурсы, связанные с клиентом S3
// (в AWS SDK v2 клиент не требует явного закрытия, но функция может быть полезна
// для обнуления ссылки или других завершающих действий)
func CloseS3(s3Client *S3) error {

	if s3Client != nil {
		// при необходимости можно обнулить встроенный клиент
		s3Client.Client = nil
		s3Client.Bucket = ""
	}

	return nil
}
