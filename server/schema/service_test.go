package schema_test

import (
	"context"
	"errors"
	"testing"

	"github.com/odpf/stencil/mocks"
	"github.com/odpf/stencil/server/domain"
	"github.com/odpf/stencil/server/schema"
	"github.com/odpf/stencil/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getSvc() (*schema.Service, *mocks.NamespaceService, *mocks.SchemaProvider, *mocks.SchemaRepository) {
	nsService := &mocks.NamespaceService{}
	schemaProvider := &mocks.SchemaProvider{}
	schemaRepo := &mocks.SchemaRepository{}
	svc := &schema.Service{
		SchemaProvider: schemaProvider,
		Repo:           schemaRepo,
		NamespaceSvc:   nsService,
	}
	return svc, nsService, schemaProvider, schemaRepo
}

func TestSchemaCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("should return error if namespace not found", func(t *testing.T) {
		svc, nsService, _, _ := getSvc()
		nsName := "testNamespace"
		nsService.On("Get", mock.Anything, nsName).Return(domain.Namespace{}, storage.NoRowsErr)
		_, err := svc.Create(ctx, nsName, "a", &domain.Metadata{}, []byte(""))
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, storage.NoRowsErr)
		nsService.AssertExpectations(t)
	})

	t.Run("should return error if schema validation fails", func(t *testing.T) {
		svc, nsService, schemaProvider, _ := getSvc()
		nsName := "testNamespace"
		data := []byte("data")
		nsService.On("Get", mock.Anything, nsName).Return(domain.Namespace{Format: "avro"}, nil)
		schemaProvider.On("ParseSchema", "protobuf", data).Return(&mocks.ParsedSchema{}, errors.New("invalid schema"))
		_, err := svc.Create(ctx, nsName, "a", &domain.Metadata{Format: "protobuf"}, data)
		assert.NotNil(t, err)
		schemaProvider.AssertExpectations(t)
	})

	t.Run("should get format from namespace if format at schema level not defined", func(t *testing.T) {
		svc, nsService, schemaProvider, _ := getSvc()
		nsName := "testNamespace"
		data := []byte("data")
		nsService.On("Get", mock.Anything, nsName).Return(domain.Namespace{Format: "protobuf"}, nil)
		schemaProvider.On("ParseSchema", "protobuf", data).Return(&mocks.ParsedSchema{}, errors.New("invalid schema"))
		_, err := svc.Create(ctx, nsName, "a", &domain.Metadata{}, data)
		assert.NotNil(t, err)
		schemaProvider.AssertExpectations(t)
	})

	t.Run("should skip compatibility check if previous latest schema not present", func(t *testing.T) {
		svc, nsService, schemaProvider, schemaRepo := getSvc()
		scFile := &domain.SchemaFile{}
		parsedSchema := &mocks.ParsedSchema{}
		nsName := "testNamespace"
		data := []byte("data")
		nsService.On("Get", mock.Anything, nsName).Return(domain.Namespace{Format: "protobuf"}, nil)
		schemaProvider.On("ParseSchema", "protobuf", data).Return(parsedSchema, nil)
		schemaRepo.On("GetLatestSchema", mock.Anything, nsName, "a").Return([]byte(""), storage.NoRowsErr)
		parsedSchema.On("GetCanonicalValue").Return(scFile)
		schemaRepo.On("CreateSchema", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(int32(1), nil)
		scInfo, err := svc.Create(ctx, nsName, "a", &domain.Metadata{}, data)
		assert.NoError(t, err)
		assert.Equal(t, scInfo.Version, int32(1))
		schemaRepo.AssertExpectations(t)
		nsService.AssertExpectations(t)
	})
	t.Run("should return error if unable to get prev latest schema", func(t *testing.T) {
		svc, nsService, schemaProvider, schemaRepo := getSvc()
		parsedSchema := &mocks.ParsedSchema{}
		nsName := "testNamespace"
		data := []byte("data")
		nsService.On("Get", mock.Anything, nsName).Return(domain.Namespace{Format: "protobuf"}, nil)
		schemaProvider.On("ParseSchema", "protobuf", data).Return(parsedSchema, nil)
		schemaRepo.On("GetLatestSchema", mock.Anything, nsName, "a").Return([]byte(""), errors.New("some other error apart from noRowsError"))
		_, err := svc.Create(ctx, nsName, "a", &domain.Metadata{}, data)
		assert.Error(t, err)
		schemaRepo.AssertExpectations(t)
		nsService.AssertExpectations(t)
	})
	t.Run("should return error if previous latest schema is valid", func(t *testing.T) {
		svc, nsService, schemaProvider, schemaRepo := getSvc()
		parsedSchema := &mocks.ParsedSchema{}
		prevParsedSchema := &mocks.ParsedSchema{}
		nsName := "testNamespace"
		data := []byte("data")
		prevData := []byte("some prev data")
		nsService.On("Get", mock.Anything, nsName).Return(domain.Namespace{Format: "protobuf", Compatibility: "BACKWARD"}, nil)
		schemaProvider.On("ParseSchema", "protobuf", data).Return(parsedSchema, nil).Once()
		schemaRepo.On("GetLatestSchema", mock.Anything, nsName, "a").Return(prevData, nil)
		schemaProvider.On("ParseSchema", "protobuf", prevData).Return(prevParsedSchema, errors.New("parse error")).Once()
		_, err := svc.Create(ctx, nsName, "a", &domain.Metadata{Compatibility: "FORWARD"}, data)
		assert.Error(t, err)
		schemaRepo.AssertExpectations(t)
		nsService.AssertExpectations(t)
		parsedSchema.AssertExpectations(t)
	})
	t.Run("should return error if compatibility check fails", func(t *testing.T) {
		svc, nsService, schemaProvider, schemaRepo := getSvc()
		parsedSchema := &mocks.ParsedSchema{}
		prevParsedSchema := &mocks.ParsedSchema{}
		nsName := "testNamespace"
		data := []byte("data")
		prevData := []byte("some prev data")
		nsService.On("Get", mock.Anything, nsName).Return(domain.Namespace{Format: "protobuf", Compatibility: "BACKWARD"}, nil)
		schemaProvider.On("ParseSchema", "protobuf", data).Return(parsedSchema, nil).Once()
		schemaRepo.On("GetLatestSchema", mock.Anything, nsName, "a").Return(prevData, nil)
		schemaProvider.On("ParseSchema", "protobuf", prevData).Return(prevParsedSchema, nil).Once()
		parsedSchema.On("IsBackwardCompatible", prevParsedSchema).Return(errors.New("compatibility error"))
		_, err := svc.Create(ctx, nsName, "a", &domain.Metadata{}, data)
		assert.Error(t, err)
		schemaRepo.AssertExpectations(t)
		nsService.AssertExpectations(t)
		parsedSchema.AssertExpectations(t)
	})
	t.Run("should return error if specified forward compatibility check fails", func(t *testing.T) {
		svc, nsService, schemaProvider, schemaRepo := getSvc()
		parsedSchema := &mocks.ParsedSchema{}
		prevParsedSchema := &mocks.ParsedSchema{}
		nsName := "testNamespace"
		data := []byte("data")
		prevData := []byte("some prev data")
		nsService.On("Get", mock.Anything, nsName).Return(domain.Namespace{Format: "protobuf", Compatibility: "BACKWARD"}, nil)
		schemaProvider.On("ParseSchema", "protobuf", data).Return(parsedSchema, nil).Once()
		schemaRepo.On("GetLatestSchema", mock.Anything, nsName, "a").Return(prevData, nil)
		schemaProvider.On("ParseSchema", "protobuf", prevData).Return(prevParsedSchema, nil).Once()
		parsedSchema.On("IsForwardCompatible", prevParsedSchema).Return(errors.New("compatibility error"))
		_, err := svc.Create(ctx, nsName, "a", &domain.Metadata{Compatibility: "FORWARD"}, data)
		assert.Error(t, err)
		schemaRepo.AssertExpectations(t)
		nsService.AssertExpectations(t)
		parsedSchema.AssertExpectations(t)
	})
}
