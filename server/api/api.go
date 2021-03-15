package api

import (
	"context"

	"github.com/odpf/stencil/server/models"
)

//StoreService Service Interface for interacting with backend store
type StoreService interface {
	ListNames(...string) []string
	ListVersions(...string) []string
	Upload(context.Context, *models.DescriptorPayload) error
	Download(context.Context, *models.FileDownload) (*models.FileData, error)
}

//API holds all handlers
type API struct {
	Store StoreService
}
