package proto

import (
	"context"

	"google.golang.org/protobuf/reflect/protoregistry"
)

// Service handles proto CRUD operations
type Service struct {
	repo *Repository
}

// GetNames returns list of available proto descriptorset names under specified namespace
func (s *Service) GetNames(ctx context.Context, namespace string) ([]string, error) {
	return s.repo.ListNames(ctx, namespace)
}

// GetVersions returns list of available versions
func (s *Service) GetVersions(ctx context.Context, namespace, name string) ([]string, error) {
	return s.repo.ListVersions(ctx, namespace, name)
}

// GetLatestVersion returns latest version number
func (s *Service) GetLatestVersion(ctx context.Context, namespace, name string) (string, error) {
	return s.repo.LatestVersion(ctx, namespace, name)
}

func (s *Service) Exists(ctx context.Context, snapshot *Snapshot) bool {
	return s.repo.Exists(ctx, snapshot)
}

// Put stores proto schema details in DB
func (s *Service) Put(ctx context.Context, snapshot *Snapshot, currentData []byte, dryRun bool) error {
	var err error
	var files *protoregistry.Files
	if files, err = getRegistry(currentData); err != nil {
		return err
	}
	dbFiles := ToProtobufDBFiles(files)
	return s.repo.Put(ctx, snapshot, dbFiles)
}

// Get returns proto schema details from DB
func (s *Service) Get(ctx context.Context, snapshot *Snapshot, names []string) (data []byte, err error) {
	dbData, err := s.repo.Get(ctx, snapshot, names)
	if err != nil {
		return
	}
	data, err = FromByteArrayToFileDescriptorSet(dbData)
	return
}

// NewService creates new instance of proto service
func NewService(r *Repository) *Service {
	return &Service{r}
}
