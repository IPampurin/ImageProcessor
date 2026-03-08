package service

import (
	"context"

	"github.com/IPampurin/ImageProcessor/pkg/manager/db"
	"github.com/IPampurin/ImageProcessor/pkg/manager/s3"
)

type Service struct {
	image  db.ImageFileMethods
	outbox db.OutboxMethods
	s3     s3.S3Methods
}

func InitService(ctx context.Context, storage *db.DataBase, storageS3 *s3.S3) *Service {

	svc := &Service{
		image:  storage,   // *db.DataBase реализует ImageFileMethods
		outbox: storage,   // *db.DataBase реализует OutboxMethods
		s3:     storageS3, // *s3.S3 реализует S3Methods
	}

	return svc
}
