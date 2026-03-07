package service

import (
	"context"

	"github.com/IPampurin/ImageProcessor/pkg/domain"
)

// ServiceMethods
type ServiceMethods interface {

	// SaveImageToDB
	SaveImageToDB(ctx context.Context, imageTask *domain.ImageTask) error

	// SaveMessageToOutbox(ctx context.Context, )
}
