package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Upload сохраняет файл в хранилище по указанному ключу (пути)
func (s *S3) Upload(ctx context.Context, bucket, path string, reader io.Reader, size int64, contentType string) error {

	input := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(path),
		Body:          reader,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(size),
	}
	_, err := s.Client.PutObject(ctx, input)

	return err
}

// Download возвращает ReadCloser для чтения файла по указанному ключу (пути)
func (s *S3) Download(ctx context.Context, bucket, path string) (io.ReadCloser, error) {

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	}
	output, err := s.Client.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}

	return output.Body, nil
}

// Delete удаляет файл по указанному ключу (пути)
func (s *S3) Delete(ctx context.Context, bucket, path string) error {

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	}
	_, err := s.Client.DeleteObject(ctx, input)

	return err
}
