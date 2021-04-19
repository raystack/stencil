package api

import (
	"context"

	"github.com/odpf/stencil/server/models"
)

//StoreService Service Interface for interacting with backend store
type StoreService interface {
	ListNames(...string) ([]string, error)
	ListVersions(...string) ([]string, error)
	Upload(context.Context, *models.DescriptorPayload) error
	Download(context.Context, *models.FileDownload) (*models.FileData, error)
	StoreMetadata(ctx context.Context, payload *models.MetadataPayload) error
	GetMetadata(ctx context.Context, payload *models.GetMetadata) (*models.MetadataFile, error)
}

//API holds all handlers
type API struct {
	Store StoreService
}
