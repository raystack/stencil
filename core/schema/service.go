package schema

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/raystack/stencil/core/namespace"
	"github.com/raystack/stencil/internal/store"
)

// NewService creates a new schema service with the given dependencies.
func NewService(repo Repository, provider Provider, nsSvc NamespaceService, cache Cache) *Service {
	return &Service{
		repo:             repo,
		provider:         provider,
		cache:            cache,
		namespaceService: nsSvc,
	}
}

type NamespaceService interface {
	Get(ctx context.Context, name string) (namespace.Namespace, error)
}

// Service provides schema management operations.
type Service struct {
	provider         Provider
	repo             Repository
	cache            Cache
	namespaceService NamespaceService
}

func (s *Service) cachedGetSchema(ctx context.Context, nsName, schemaName string, version int32) ([]byte, error) {
	key := schemaKeyFunc(nsName, schemaName, version)
	val, found := s.cache.Get(key)
	if !found {
		var data []byte
		var err error
		data, err = s.repo.Get(ctx, nsName, schemaName, version)
		if err != nil {
			return data, err
		}
		s.cache.Set(key, data, int64(len(data)))
		return data, err
	}
	return getBytes(val), nil
}

func (s *Service) CheckCompatibility(ctx context.Context, nsName, schemaName, compatibility string, data []byte) error {
	ns, err := s.namespaceService.Get(ctx, nsName)
	if err != nil {
		return err
	}
	compatibility = getNonEmpty(compatibility, ns.Compatibility)
	parsedSchema, err := s.provider.ParseSchema(ns.Format, data)
	if err != nil {
		return err
	}
	return s.checkCompatibility(ctx, nsName, schemaName, ns.Format, compatibility, parsedSchema)
}

func (s *Service) checkCompatibility(ctx context.Context, nsName, schemaName, format, compatibility string, current ParsedSchema) error {
	prevMeta, prevSchemaData, err := s.GetLatest(ctx, nsName, schemaName)
	if err != nil {
		if errors.Is(err, store.NoRowsErr) {
			return nil
		}
		return err
	}
	prevSchema, err := s.provider.ParseSchema(prevMeta.Format, prevSchemaData)
	if err != nil {
		return err
	}
	checkerFn := getCompatibilityChecker(compatibility)
	return checkerFn(current, []ParsedSchema{prevSchema})
}

// Create validates, parses, and stores a new schema version.
func (s *Service) Create(ctx context.Context, nsName string, schemaName string, metadata *Metadata, data []byte) (SchemaInfo, error) {
	var scInfo SchemaInfo
	ns, err := s.namespaceService.Get(ctx, nsName)
	if err != nil {
		return scInfo, err
	}
	format := getNonEmpty(metadata.Format, ns.Format)
	compatibility := getNonEmpty(metadata.Compatibility, ns.Compatibility)
	parsedSchema, err := s.provider.ParseSchema(format, data)
	if err != nil {
		return scInfo, err
	}
	if err := s.checkCompatibility(ctx, nsName, schemaName, format, compatibility, parsedSchema); err != nil {
		return scInfo, err
	}
	sf := parsedSchema.GetCanonicalValue()
	mergedMetadata := &Metadata{
		Format:        format,
		Compatibility: compatibility,
	}
	versionID := getIDforSchema(nsName, schemaName, sf.ID)
	version, err := s.repo.Create(ctx, nsName, schemaName, mergedMetadata, versionID, sf)
	return SchemaInfo{
		Version:  version,
		ID:       versionID,
		Location: fmt.Sprintf("/v1beta1/namespaces/%s/schemas/%s/versions/%d", nsName, schemaName, version),
	}, err
}

func (s *Service) withMetadata(ctx context.Context, namespace, schemaName string, getData func() ([]byte, error)) (*Metadata, []byte, error) {
	var data []byte
	meta, err := s.repo.GetMetadata(ctx, namespace, schemaName)
	if err != nil {
		return meta, data, err
	}
	data, err = getData()
	return meta, data, err
}

// Get retrieves a specific schema version by namespace, name, and version number.
func (s *Service) Get(ctx context.Context, namespace string, schemaName string, version int32) (*Metadata, []byte, error) {
	return s.withMetadata(ctx, namespace, schemaName, func() ([]byte, error) { return s.cachedGetSchema(ctx, namespace, schemaName, version) })
}

// Delete removes a schema and all its versions, and invalidates cached entries.
func (s *Service) Delete(ctx context.Context, namespace string, schemaName string) error {
	versions, err := s.repo.ListVersions(ctx, namespace, schemaName)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, namespace, schemaName); err != nil {
		return err
	}
	for _, v := range versions {
		s.cache.Del(schemaKeyFunc(namespace, schemaName, v))
	}
	return nil
}

// DeleteVersion removes a specific version of a schema and invalidates the cached entry.
func (s *Service) DeleteVersion(ctx context.Context, namespace string, schemaName string, version int32) error {
	err := s.repo.DeleteVersion(ctx, namespace, schemaName, version)
	if err == nil {
		s.cache.Del(schemaKeyFunc(namespace, schemaName, version))
	}
	return err
}

// GetLatest retrieves the latest version of a schema.
func (s *Service) GetLatest(ctx context.Context, namespace string, schemaName string) (*Metadata, []byte, error) {
	version, err := s.repo.GetLatestVersion(ctx, namespace, schemaName)
	if err != nil {
		return nil, nil, err
	}
	return s.Get(ctx, namespace, schemaName, version)
}

// GetMetadata retrieves the metadata for a schema.
func (s *Service) GetMetadata(ctx context.Context, namespace, schemaName string) (*Metadata, error) {
	return s.repo.GetMetadata(ctx, namespace, schemaName)
}

// UpdateMetadata updates the metadata for a schema.
func (s *Service) UpdateMetadata(ctx context.Context, namespace, schemaName string, meta *Metadata) (*Metadata, error) {
	return s.repo.UpdateMetadata(ctx, namespace, schemaName, meta)
}

// List returns all schemas in a namespace.
func (s *Service) List(ctx context.Context, namespaceID string) ([]Schema, error) {
	return s.repo.List(ctx, namespaceID)
}

// ListVersions returns all version numbers for a schema.
func (s *Service) ListVersions(ctx context.Context, namespaceID string, schemaName string) ([]int32, error) {
	return s.repo.ListVersions(ctx, namespaceID, schemaName)
}

func getIDforSchema(ns, schema, dataUUID string) string {
	key := fmt.Sprintf("%s-%s-%s", ns, schema, dataUUID)
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(key)).String()
}
