package snapshot

import (
	"context"

	"github.com/odpf/stencil/models"
	"github.com/odpf/stencil/storage"
)

// Service handles proto CRUD operations
type Service struct {
	store storage.Store
}

func (s *Service) Exists(ctx context.Context, snapshot *models.Snapshot) bool {
	return s.store.ExistsSnapshot(ctx, snapshot)
}

func (s *Service) List(ctx context.Context, snapshot *models.Snapshot) ([]*models.Snapshot, error) {
	return s.store.ListSnapshots(ctx, snapshot)
}

func (s *Service) GetSnapshotByFields(ctx context.Context, namespace, name, version string, latest bool) (*models.Snapshot, error) {
	return s.store.GetSnapshotByFields(ctx, namespace, name, version, latest)
}

func (s *Service) GetSnapshotByID(ctx context.Context, id int64) (*models.Snapshot, error) {
	return s.store.GetSnapshotByID(ctx, id)
}

func (s *Service) UpdateLatestVersion(ctx context.Context, snapshot *models.Snapshot) error {
	return s.store.UpdateSnapshotLatestVersion(ctx, snapshot)
}

// NewService creates new instance of proto service
func NewService(store storage.Store) *Service {
	return &Service{store: store}
}
