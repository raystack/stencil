package proto

import (
	"context"

	"github.com/odpf/stencil/models"
	"github.com/odpf/stencil/storage"
)

// Service handles proto CRUD operations
type Service struct {
	store storage.Store
}

// Validate checks if current data is backward compatible against previous stable data
func (s *Service) Validate(ctx context.Context, cs *models.Snapshot, data []byte, rulesToSkip []string) error {
	var err error
	prevSt, err := s.store.GetSnapshotByFields(ctx, cs.Namespace, cs.Name, "", nil)
	if err == models.ErrSnapshotNotFound {
		return nil
	}
	// no need to handle error here. Since without snapshot, data won't exist.
	// If snapshot exist and data is nil, then validation still passes as it's treated as completely new
	prevData, _ := s.Get(ctx, prevSt, []string{})
	return Compare(data, prevData, rulesToSkip)
}

func (s *Service) Merge(ctx context.Context, prevData, data []byte) ([]byte, error) {
	return Merge(data, prevData)
}

// Insert stores proto schema details in DB after backward compatible check succeeds
func (s *Service) Insert(ctx context.Context, snapshot *models.Snapshot, data []byte) error {
	files, _ := getRegistry(data)
	dbFiles := toProtobufDBFiles(files)
	err := s.store.PutSchema(ctx, snapshot, dbFiles)
	if err != nil {
		return err
	}
	return s.store.UpdateSnapshotLatestVersion(ctx, snapshot)
}

// Get returns proto schema details from DB
func (s *Service) Get(ctx context.Context, snapshot *models.Snapshot, names []string) (data []byte, err error) {
	dbData, err := s.store.GetSchema(ctx, snapshot, names)
	if err != nil {
		return
	}
	data, err = fromByteArrayToFileDescriptorSet(dbData)
	return
}

// NewService creates new instance of proto service
func NewService(store storage.Store) *Service {
	return &Service{store: store}
}
