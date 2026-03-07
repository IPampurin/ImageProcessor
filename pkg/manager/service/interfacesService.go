package service

import (
	"context"

	"github.com/IPampurin/ImageProcessor/pkg/manager/api"
	"github.com/google/uuid"
)

// ServiceMethods
type ServiceMethods interface {

	// UploadImage
	UploadImage(ctx context.Context, req *api.UploadRequest) (uuid.UUID, error)

	// SaveMessageToOutbox(ctx context.Context, )
}
