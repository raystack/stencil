package api

import (
	"context"

	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"github.com/odpf/stencil/server/schema"
)

func (a *API) CreateSchema(ctx context.Context, in *stencilv1.CreateSchemaRequest) (*stencilv1.CreateSchemaResponse, error) {
	metadata := &schema.Metadata{Format: in.GetFormat().String(), Compatibility: in.GetCompatibility().String()}
	sc, err := a.Schema.Create(ctx, in.NamespaceId, in.SchemaId, metadata, in.GetData())
	return &stencilv1.CreateSchemaResponse{
		Version:  sc.Version,
		Id:       sc.ID,
		Location: sc.Location,
	}, err
}

func (a *API) ListSchemas(ctx context.Context, in *stencilv1.ListSchemasRequest) (*stencilv1.ListSchemasResponse, error) {
	schemas, err := a.Schema.List(ctx, in.Id)
	return &stencilv1.ListSchemasResponse{Schemas: schemas}, err
}

func (a *API) GetLatestSchema(ctx context.Context, in *stencilv1.GetLatestSchemaRequest) (*stencilv1.GetLatestSchemaResponse, error) {
	data, err := a.Schema.GetLatest(ctx, in.NamespaceId, in.SchemaId)
	return &stencilv1.GetLatestSchemaResponse{
		Data: data,
	}, err
}

func (a *API) GetSchema(ctx context.Context, in *stencilv1.GetSchemaRequest) (*stencilv1.GetSchemaResponse, error) {
	data, err := a.Schema.Get(ctx, in.NamespaceId, in.SchemaId, in.GetVersionId())
	return &stencilv1.GetSchemaResponse{
		Data: data,
	}, err
}

func (a *API) ListVersions(ctx context.Context, in *stencilv1.ListVersionsRequest) (*stencilv1.ListVersionsResponse, error) {
	versions, err := a.Schema.ListVersions(ctx, in.NamespaceId, in.SchemaId)
	return &stencilv1.ListVersionsResponse{Versions: versions}, err
}

func (a *API) GetSchemaMetadata(ctx context.Context, in *stencilv1.GetSchemaMetadataRequest) (*stencilv1.GetSchemaMetadataResponse, error) {
	meta, err := a.Schema.GetMetadata(ctx, in.NamespaceId, in.SchemaId)
	return &stencilv1.GetSchemaMetadataResponse{
		Format:        stencilv1.Schema_Format(stencilv1.Schema_Format_value[meta.Format]),
		Compatibility: stencilv1.Schema_Compatibility(stencilv1.Schema_Compatibility_value[meta.Compatibility]),
		Authority:     meta.Authority,
	}, err
}

func (a *API) UpdateSchemaMetadata(ctx context.Context, in *stencilv1.UpdateSchemaMetadataRequest) (*stencilv1.UpdateSchemaMetadataResponse, error) {
	meta, err := a.Schema.UpdateMetadata(ctx, in.NamespaceId, in.SchemaId, &schema.Metadata{
		Compatibility: in.Compatibility.String(),
	})
	return &stencilv1.UpdateSchemaMetadataResponse{
		Format:        stencilv1.Schema_Format(stencilv1.Schema_Format_value[meta.Format]),
		Compatibility: stencilv1.Schema_Compatibility(stencilv1.Schema_Compatibility_value[meta.Compatibility]),
		Authority:     meta.Authority,
	}, err
}

func (a *API) DeleteSchema(ctx context.Context, in *stencilv1.DeleteSchemaRequest) (*stencilv1.DeleteSchemaResponse, error) {
	err := a.Schema.Delete(ctx, in.NamespaceId, in.SchemaId)
	message := "success"
	if err != nil {
		message = "failed"
	}
	return &stencilv1.DeleteSchemaResponse{
		Message: message,
	}, err
}

func (a *API) DeleteVersion(ctx context.Context, in *stencilv1.DeleteVersionRequest) (*stencilv1.DeleteVersionResponse, error) {
	err := a.Schema.DeleteVersion(ctx, in.NamespaceId, in.SchemaId, in.GetVersionId())
	message := "success"
	if err != nil {
		message = "failed"
	}
	return &stencilv1.DeleteVersionResponse{
		Message: message,
	}, err
}
