package api

import (
	"context"

	"github.com/odpf/stencil/server/api/v1/pb"
	"github.com/odpf/stencil/server/snapshot"
	"google.golang.org/grpc/health/grpc_health_v1"
)

//StoreService Service Interface for storage and validation
type StoreService interface {
	Validate(context.Context, *snapshot.Snapshot, []byte, []string) error
	Insert(context.Context, *snapshot.Snapshot, []byte) error
	Get(context.Context, *snapshot.Snapshot, []string) ([]byte, error)
}

// MetadataService Service Interface for metadata store
type MetadataService interface {
	Exists(context.Context, *snapshot.Snapshot) bool
	List(context.Context, *snapshot.Snapshot) ([]*snapshot.Snapshot, error)
	GetSnapshotByFields(context.Context, string, string, string, bool) (*snapshot.Snapshot, error)
	GetSnapshotByID(context.Context, int64) (*snapshot.Snapshot, error)
	UpdateLatestVersion(context.Context, *snapshot.Snapshot) error
}

//API holds all handlers
type API struct {
	pb.UnimplementedStencilServiceServer
	grpc_health_v1.UnimplementedHealthServer
	Store    StoreService
	Metadata MetadataService
}
