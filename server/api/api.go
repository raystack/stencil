package api

import (
	"context"

	"github.com/odpf/stencil/models"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"google.golang.org/grpc/health/grpc_health_v1"
)

//StoreService Service Interface for storage and validation
type StoreService interface {
	Validate(context.Context, *models.Snapshot, []byte, []string) error
	Insert(context.Context, *models.Snapshot, []byte) error
	Get(context.Context, *models.Snapshot, []string) ([]byte, error)
}

// MetadataService Service Interface for metadata store
type MetadataService interface {
	Exists(context.Context, *models.Snapshot) bool
	List(context.Context, *models.Snapshot) ([]*models.Snapshot, error)
	GetSnapshotByFields(context.Context, string, string, string, bool) (*models.Snapshot, error)
	GetSnapshotByID(context.Context, int64) (*models.Snapshot, error)
	UpdateLatestVersion(context.Context, *models.Snapshot) error
}

//API holds all handlers
type API struct {
	stencilv1.UnimplementedStencilServiceServer
	grpc_health_v1.UnimplementedHealthServer
	Store    StoreService
	Metadata MetadataService
}
