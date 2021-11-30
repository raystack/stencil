package schema

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/odpf/stencil/server/domain"
	"github.com/odpf/stencil/storage"
)

func getNonEmpty(args ...string) string {
	for _, a := range args {
		if a != "" {
			return a
		}
	}
	return ""
}

type Service struct {
	SchemaProvider SchemaProvider
	Repo           domain.SchemaRepository
	NamespaceSvc   domain.NamespaceService
}

func (s *Service) CheckCompatibility(ctx context.Context, nsName, schemaName, format, compatibility string, current ParsedSchema) error {
	prevSchemaData, err := s.GetLatest(ctx, nsName, schemaName)
	if err != nil {
		if errors.Is(err, storage.NoRowsErr) {
			return nil
		}
		return err
	}
	prevSchema, err := s.SchemaProvider.ParseSchema(format, prevSchemaData)
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
	if err := s.CheckCompatibility(ctx, nsName, schemaName, format, compatibility, parsedSchema); err != nil {
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
		Location: fmt.Sprintf("/v1/namespaces/%s/schemas/%s/versions/%d", nsName, schemaName, version),
	}, err
}

func (s *Service) Get(ctx context.Context, namespace string, schemaName string, version int32) ([]byte, error) {
	return s.Repo.GetSchema(ctx, namespace, schemaName, version)
}

func (s *Service) Delete(ctx context.Context, namespace string, schemaName string) error {
	return s.Repo.DeleteSchema(ctx, namespace, schemaName)
}

func (s *Service) DeleteVersion(ctx context.Context, namespace string, schemaName string, version int32) error {
	return s.Repo.DeleteVersion(ctx, namespace, schemaName, version)
}

func (s *Service) GetLatest(ctx context.Context, namespace string, schemaName string) ([]byte, error) {
	return s.Repo.GetLatestSchema(ctx, namespace, schemaName)
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
