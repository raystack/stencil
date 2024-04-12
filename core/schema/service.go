package schema

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/goto/stencil/config"
	"github.com/goto/stencil/core/changedetector"
	"github.com/goto/stencil/core/namespace"
	"github.com/goto/stencil/internal/store"
	"github.com/goto/stencil/pkg/newrelic"
)

func NewService(repo Repository, provider Provider, nsSvc NamespaceService,
	cache Cache, nr newrelic.Service, cds ChangeDetectorService,
	producer Producer, config *config.Config) *Service {
	return &Service{
		repo:                  repo,
		provider:              provider,
		cache:                 cache,
		namespaceService:      nsSvc,
		newrelic:              nr,
		changeDetectorService: cds,
		producer:              producer,
		config:                config,
	}
}

type NamespaceService interface {
	Get(ctx context.Context, name string) (namespace.Namespace, error)
}

type Service struct {
	provider              Provider
	repo                  Repository
	cache                 Cache
	namespaceService      NamespaceService
	newrelic              newrelic.Service
	changeDetectorService ChangeDetectorService
	producer              Producer
	config                *config.Config
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
	endFunc := s.newrelic.StartGenericSegment(ctx, "Compatibility checker")
	defer endFunc()
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

func (s *Service) Create(ctx context.Context, nsName string, schemaName string, metadata *Metadata, data []byte) (SchemaInfo, error) {
	endFunc := s.newrelic.StartGenericSegment(ctx, "Create Schema Info")
	defer endFunc()
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
	_, prevSchemaData, _ := s.GetLatest(ctx, nsName, schemaName)
	version, err := s.repo.Create(ctx, nsName, schemaName, mergedMetadata, versionID, sf)
	changeRequest := &changedetector.ChangeRequest{
		NamespaceName: nsName,
		SchemaName:    schemaName,
		Version:       version,
		OldData:       prevSchemaData,
		NewData:       data,
	}
	go s.identifySchemaChange(ctx, changeRequest)
	return SchemaInfo{
		Version:  version,
		ID:       versionID,
		Location: fmt.Sprintf("/v1beta1/namespaces/%s/schemas/%s/versions/%d", nsName, schemaName, version),
	}, err
}

func (s *Service) identifySchemaChange(ctx context.Context, request *changedetector.ChangeRequest) error {
	endFunc := s.newrelic.StartGenericSegment(ctx, "Identify Schema Change")
	defer endFunc()
	sce, err := s.changeDetectorService.IdentifySchemaChange(ctx, request)
	if err != nil {
		log.Printf("got error while identifying schema change for namespace : %s, schema: %s, version: %d, %v", request.NamespaceName, request.SchemaName, request.Version, err)
		return err
	}
	log.Printf("schema change result %s", sce.String())
	schemaChangeTopic := s.config.SchemaChange.KafkaTopic
	retryInterval := time.Duration(s.config.KafkaProducer.RetryInterval) * time.Millisecond
	if err := s.producer.PushMessagesWithRetries(schemaChangeTopic, sce, s.config.KafkaProducer.Retries, retryInterval); err != nil {
		log.Printf("unable to push message to Kafka topic %s for schema change event %s: %s", schemaChangeTopic, sce, err.Error())
		return err
	}
	log.Printf("successfully pushed message to kafka topic %s", schemaChangeTopic)
	return nil
}

func (s *Service) withMetadata(ctx context.Context, namespace, schemaName string, getData func() ([]byte, error)) (*Metadata, []byte, error) {
	endFunc := s.newrelic.StartGenericSegment(ctx, "GetMetaData")
	defer endFunc()
	var data []byte
	meta, err := s.repo.GetMetadata(ctx, namespace, schemaName)
	if err != nil {
		return meta, data, err
	}

	dataSegmentEndFunc := s.newrelic.StartGenericSegment(ctx, "GetData")
	data, err = getData()
	dataSegmentEndFunc()
	return meta, data, err
}

func (s *Service) Get(ctx context.Context, namespace string, schemaName string, version int32) (*Metadata, []byte, error) {
	return s.withMetadata(ctx, namespace, schemaName, func() ([]byte, error) { return s.cachedGetSchema(ctx, namespace, schemaName, version) })
}

func (s *Service) Delete(ctx context.Context, namespace string, schemaName string) error {
	return s.repo.Delete(ctx, namespace, schemaName)
}

func (s *Service) DeleteVersion(ctx context.Context, namespace string, schemaName string, version int32) error {
	return s.repo.DeleteVersion(ctx, namespace, schemaName, version)
}

func (s *Service) GetLatest(ctx context.Context, namespace string, schemaName string) (*Metadata, []byte, error) {
	version, err := s.repo.GetLatestVersion(ctx, namespace, schemaName)
	if err != nil {
		return nil, nil, err
	}
	return s.Get(ctx, namespace, schemaName, version)
}

func (s *Service) GetMetadata(ctx context.Context, namespace, schemaName string) (*Metadata, error) {
	return s.repo.GetMetadata(ctx, namespace, schemaName)
}

func (s *Service) UpdateMetadata(ctx context.Context, namespace, schemaName string, meta *Metadata) (*Metadata, error) {
	return s.repo.UpdateMetadata(ctx, namespace, schemaName, meta)
}

func (s *Service) List(ctx context.Context, namespaceID string) ([]Schema, error) {
	return s.repo.List(ctx, namespaceID)
}

func (s *Service) ListVersions(ctx context.Context, namespaceID string, schemaName string) ([]int32, error) {
	return s.repo.ListVersions(ctx, namespaceID, schemaName)
}

func getIDforSchema(ns, schema, dataUUID string) string {
	key := fmt.Sprintf("%s-%s-%s", ns, schema, dataUUID)
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(key)).String()
}
