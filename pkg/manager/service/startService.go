package service

import (
	"context"

	"github.com/IPampurin/ImageProcessor/pkg/manager/db"
)

type Service struct {
	image  db.ImageFileMethods
	outbox db.OutboxMethods
}

func InitService(ctx context.Context, storage *db.DataBase) *Service {

	svc := &Service{
		image:  storage, // *db.DataBase реализует ImageFileMethods
		outbox: storage, // *db.DataBase реализует OutboxMethods
	}

	return svc
}
