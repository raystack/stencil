package snapshot

import (
	"context"
)

// Service handles snapshot CRUD operations
type Service struct {
	repo *Repository
}

// ListNames returns list of available proto descriptorset names under specified namespace
func (s *Service) ListNames(ctx context.Context, namespace string) ([]string, error) {
	return s.repo.ListNames(ctx, namespace)
}

// ListVersions returns list of available versions
func (s *Service) ListVersions(ctx context.Context, namespace, name string) ([]string, error) {
	return s.repo.ListVersions(ctx, namespace, name)
}

// GetSnapshot returns latest version number
func (s *Service) GetSnapshot(ctx context.Context, namespace, name, version string, latest bool) (*Snapshot, error) {
	return s.repo.GetSnapshot(ctx, namespace, name, version, latest)
}

// UpdateLatestVersion updates latest version number for snapshot
func (s *Service) UpdateLatestVersion(ctx context.Context, st *Snapshot) error {
	snapshotWithID, err := s.repo.GetSnapshot(ctx, st.Namespace, st.Name, st.Version, st.Latest)
	if err != nil {
		return err
	}
	return s.repo.UpdateLatestVersion(ctx, snapshotWithID)
}

// Exists check if snapshot exists or not
func (s *Service) Exists(ctx context.Context, snapshot *Snapshot) bool {
	return s.repo.Exists(ctx, snapshot)
}

// NewSnapshotService creates new instance of proto service
func NewSnapshotService(r *Repository) *Service {
	return &Service{r}
}
