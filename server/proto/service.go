package proto

import (
	"context"

	"github.com/odpf/stencil/server/snapshot"
)

// Service handles proto CRUD operations
type Service struct {
	repo         *Repository
	snapshotRepo *snapshot.Repository
}

// Validate checks if current data is backward compatible against previous stable data
func (s *Service) Validate(ctx context.Context, cs *snapshot.Snapshot, data []byte, rulesToSkip []string) error {
	var err error
	prevSt, err := s.snapshotRepo.GetSnapshotByFields(ctx, cs.Namespace, cs.Name, "", true)
	if err == snapshot.ErrNotFound {
		return nil
	}
	// no need to handle error here. Since without snapshot data won't exist.
	// If snapshot exist and data is nil, then validation still passes as it's treated as completely new
	prevData, _ := s.Get(ctx, prevSt, []string{})
	return Compare(data, prevData, rulesToSkip)
}

// Insert stores proto schema details in DB after backward compatible check succeeds
func (s *Service) Insert(ctx context.Context, snapshot *snapshot.Snapshot, data []byte) error {
	files, _ := getRegistry(data)
	dbFiles := toProtobufDBFiles(files)
	err := s.repo.Put(ctx, snapshot, dbFiles)
	if err != nil {
		return err
	}
	return s.snapshotRepo.UpdateLatestVersion(ctx, snapshot)
}

// Get returns proto schema details from DB
func (s *Service) Get(ctx context.Context, snapshot *snapshot.Snapshot, names []string) (data []byte, err error) {
	dbData, err := s.repo.Get(ctx, snapshot, names)
	if err != nil {
		return
	}
	data, err = fromByteArrayToFileDescriptorSet(dbData)
	return
}

// NewService creates new instance of proto service
func NewService(r *Repository, sr *snapshot.Repository) *Service {
	return &Service{repo: r, snapshotRepo: sr}
}
