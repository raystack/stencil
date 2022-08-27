package schema

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/odpf/stencil/domain"
	"github.com/odpf/stencil/internal/store"
)

func getNonEmpty(args ...string) string {
	for _, a := range args {
		if a != "" {
			return a
		}
	}
	return ""
}

func schemaKeyFunc(nsName, schema string, version int32) string {
	return fmt.Sprintf("%s-%s-%d", nsName, schema, version)
}

func getBytes(key interface{}) []byte {
	buf, _ := key.([]byte)
	return buf
}

func NewService(repo domain.SchemaRepository, provider SchemaProvider, nsSvc domain.NamespaceService, cache schemaCache) *Service {
	return &Service{
		Repo:           repo,
		SchemaProvider: provider,
		NamespaceSvc:   nsSvc,
		cache:          cache,
	}
}

type schemaCache interface {
	Get(interface{}) (interface{}, bool)
	Set(interface{}, interface{}, int64) bool
}

type Service struct {
	SchemaProvider SchemaProvider
	Repo           domain.SchemaRepository
	NamespaceSvc   domain.NamespaceService
	cache          schemaCache
}

func (s *Service) cachedGetSchema(ctx context.Context, nsName, schemaName string, version int32) ([]byte, error) {
	key := schemaKeyFunc(nsName, schemaName, version)
	val, found := s.cache.Get(key)
	if !found {
		var data []byte
		var err error
		data, err = s.Repo.GetSchema(ctx, nsName, schemaName, version)
		if err != nil {
			return data, err
		}
		s.cache.Set(key, data, int64(len(data)))
		return data, err
	}
	return getBytes(val), nil
}

func (s *Service) CheckCompatibility(ctx context.Context, nsName, schemaName, compatibility string, data []byte) error {
	ns, err := s.NamespaceSvc.Get(ctx, nsName)
	if err != nil {
		return err
	}
	compatibility = getNonEmpty(compatibility, ns.Compatibility)
	parsedSchema, err := s.SchemaProvider.ParseSchema(ns.Format, data)
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
	prevSchema, err := s.SchemaProvider.ParseSchema(prevMeta.Format, prevSchemaData)
	if err != nil {
		return err
	}
	checkerFn := getCompatibilityChecker(compatibility)
	return checkerFn(current, []ParsedSchema{prevSchema})
}

func (s *Service) Create(ctx context.Context, nsName string, schemaName string, metadata *domain.Metadata, data []byte) (domain.SchemaInfo, error) {
	var scInfo domain.SchemaInfo
	ns, err := s.NamespaceSvc.Get(ctx, nsName)
	if err != nil {
		return scInfo, err
	}
	format := getNonEmpty(metadata.Format, ns.Format)
	compatibility := getNonEmpty(metadata.Compatibility, ns.Compatibility)
	parsedSchema, err := s.SchemaProvider.ParseSchema(format, data)
	if err != nil {
		return scInfo, err
	}
	if err := s.checkCompatibility(ctx, nsName, schemaName, format, compatibility, parsedSchema); err != nil {
		return scInfo, err
	}
	sf := parsedSchema.GetCanonicalValue()
	mergedMetadata := &domain.Metadata{
		Format:        format,
		Compatibility: compatibility,
	}
	versionID := getIDforSchema(nsName, schemaName, sf.ID)
	version, err := s.Repo.CreateSchema(ctx, nsName, schemaName, mergedMetadata, versionID, sf)
	return domain.SchemaInfo{
		Version:  version,
		ID:       versionID,
		Location: fmt.Sprintf("/v1beta1/namespaces/%s/schemas/%s/versions/%d", nsName, schemaName, version),
	}, err
}

func (s *Service) withMetadata(ctx context.Context, namespace, schemaName string, getData func() ([]byte, error)) (*domain.Metadata, []byte, error) {
	var data []byte
	meta, err := s.Repo.GetSchemaMetadata(ctx, namespace, schemaName)
	if err != nil {
		return meta, data, err
	}
	data, err = getData()
	return meta, data, err
}

func (s *Service) Get(ctx context.Context, namespace string, schemaName string, version int32) (*domain.Metadata, []byte, error) {
	return s.withMetadata(ctx, namespace, schemaName, func() ([]byte, error) { return s.cachedGetSchema(ctx, namespace, schemaName, version) })
}

func (s *Service) Delete(ctx context.Context, namespace string, schemaName string) error {
	return s.Repo.DeleteSchema(ctx, namespace, schemaName)
}

func (s *Service) DeleteVersion(ctx context.Context, namespace string, schemaName string, version int32) error {
	return s.Repo.DeleteVersion(ctx, namespace, schemaName, version)
}

func (s *Service) GetLatest(ctx context.Context, namespace string, schemaName string) (*domain.Metadata, []byte, error) {
	version, err := s.Repo.GetLatestVersion(ctx, namespace, schemaName)
	if err != nil {
		return nil, nil, err
	}
	return s.Get(ctx, namespace, schemaName, version)
}

func (s *Service) GetMetadata(ctx context.Context, namespace, schemaName string) (*domain.Metadata, error) {
	return s.Repo.GetSchemaMetadata(ctx, namespace, schemaName)
}

func (s *Service) UpdateMetadata(ctx context.Context, namespace, schemaName string, meta *domain.Metadata) (*domain.Metadata, error) {
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
