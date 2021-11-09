package schema

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/odpf/stencil/server/namespace"
)

// Metadata model
type Metadata struct {
	Authority     string
	Format        string
	Compatibility string
}

type SchemaInfo struct {
	ID       string
	Version  int32
	Location string
}

type SchemaFile struct {
	ID           string
	Dependencies []string
	Types        []string
	Fields       []string
	Data         []byte
}

// Repository for Schema
type Repository interface {
	CreateSchema(ctx context.Context, namespace string, schema string, metadata *Metadata, versionID string, schemaFile *SchemaFile) (version int32, err error)
	ListSchemas(context.Context, string) ([]string, error)
	ListVersions(context.Context, string, string) ([]int32, error)
	GetSchema(context.Context, string, string, int32) ([]byte, error)
	GetLatestSchema(context.Context, string, string) ([]byte, error)
	GetSchemaMetadata(context.Context, string, string) (*Metadata, error)
	UpdateSchemaMetadata(context.Context, string, string, *Metadata) (*Metadata, error)
	DeleteVersion(context.Context, string, string, int32) error
}

type SchemaProvider interface {
	GetSchemaFile(format string, data []byte) (*SchemaFile, error)
}

type Service struct {
	SchemaProvider SchemaProvider
	Repo           Repository
	NamespaceSvc   namespace.Service
}

func getNonEmpty(args ...string) string {
	for _, a := range args {
		if a != "" {
			return a
		}
	}
	return ""
}

func (s *Service) Create(ctx context.Context, nsName string, schemaName string, metadata *Metadata, data []byte) (SchemaInfo, error) {
	var scInfo SchemaInfo
	ns, err := s.NamespaceSvc.Get(ctx, nsName)
	if err != nil {
		return scInfo, err
	}
	sf, err := s.SchemaProvider.GetSchemaFile(ns.Format, data)
	if err != nil {
		return scInfo, err
	}
	mergedMetadata := &Metadata{
		Format:        getNonEmpty(metadata.Format, ns.Format),
		Compatibility: getNonEmpty(metadata.Compatibility, ns.Compatibility),
	}
	versionID := getIDforSchema(nsName, schemaName, sf.ID)
	version, err := s.Repo.CreateSchema(ctx, nsName, schemaName, mergedMetadata, versionID, sf)
	return SchemaInfo{
		Version:  version,
		ID:       versionID,
		Location: fmt.Sprintf("/v1/namespaces/%s/schemas/%s/versions/%d", nsName, schemaName, version),
	}, err
}

func (s *Service) Get(ctx context.Context, namespace string, schemaName string, version int32) ([]byte, error) {
	return s.Repo.GetSchema(ctx, namespace, schemaName, version)
}

func (s *Service) DeleteVersion(ctx context.Context, namespace string, schemaName string, version int32) error {
	return s.Repo.DeleteVersion(ctx, namespace, schemaName, version)
}

func (s *Service) GetLatest(ctx context.Context, namespace string, schemaName string) ([]byte, error) {
	return s.Repo.GetLatestSchema(ctx, namespace, schemaName)
}

func (s *Service) GetMetadata(ctx context.Context, namespace, schemaName string) (*Metadata, error) {
	return s.Repo.GetSchemaMetadata(ctx, namespace, schemaName)
}

func (s *Service) UpdateMetadata(ctx context.Context, namespace, schemaName string, meta *Metadata) (*Metadata, error) {
	return s.Repo.UpdateSchemaMetadata(ctx, namespace, schemaName, meta)
}

func (s *Service) List(ctx context.Context, namespaceID string) ([]string, error) {
	return s.Repo.ListSchemas(ctx, namespaceID)
}

func (s *Service) ListVersions(ctx context.Context, namespaceID string, schemaName string) ([]int32, error) {
	return s.Repo.ListVersions(ctx, namespaceID, schemaName)
}

func getIDforSchema(ns, schema, dataUUID string) string {
	key := fmt.Sprintf("%s-%s-%s", ns, schema, dataUUID)
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(key)).String()
}
