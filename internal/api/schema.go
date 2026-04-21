package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/raystack/stencil/core/schema"
	stencilv1beta1 "github.com/raystack/stencil/gen/raystack/stencil/v1beta1"
)

func schemaToProto(s schema.Schema) *stencilv1beta1.Schema {
	return &stencilv1beta1.Schema{
		Name:          s.Name,
		Format:        stencilv1beta1.Schema_Format(stencilv1beta1.Schema_Format_value[s.Format]),
		Compatibility: stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[s.Compatibility]),
		Authority:     s.Authority,
	}
}

func (a *API) CreateSchema(ctx context.Context, req *connect.Request[stencilv1beta1.CreateSchemaRequest]) (*connect.Response[stencilv1beta1.CreateSchemaResponse], error) {
	metadata := &schema.Metadata{Format: req.Msg.GetFormat().String(), Compatibility: req.Msg.GetCompatibility().String()}
	sc, err := a.schema.Create(ctx, req.Msg.GetNamespaceId(), req.Msg.GetSchemaId(), metadata, req.Msg.GetData())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&stencilv1beta1.CreateSchemaResponse{
		Version:  sc.Version,
		Id:       sc.ID,
		Location: sc.Location,
	}), nil
}

func (a *API) CheckCompatibility(ctx context.Context, req *connect.Request[stencilv1beta1.CheckCompatibilityRequest]) (*connect.Response[stencilv1beta1.CheckCompatibilityResponse], error) {
	err := a.schema.CheckCompatibility(ctx, req.Msg.GetNamespaceId(), req.Msg.GetSchemaId(), req.Msg.GetCompatibility().String(), req.Msg.GetData())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&stencilv1beta1.CheckCompatibilityResponse{}), nil
}

func (a *API) ListSchemas(ctx context.Context, req *connect.Request[stencilv1beta1.ListSchemasRequest]) (*connect.Response[stencilv1beta1.ListSchemasResponse], error) {
	schemas, err := a.schema.List(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	var ss []*stencilv1beta1.Schema
	for _, s := range schemas {
		ss = append(ss, schemaToProto(s))
	}
	return connect.NewResponse(&stencilv1beta1.ListSchemasResponse{Schemas: ss}), nil
}

func (a *API) GetLatestSchema(ctx context.Context, req *connect.Request[stencilv1beta1.GetLatestSchemaRequest]) (*connect.Response[stencilv1beta1.GetLatestSchemaResponse], error) {
	_, data, err := a.schema.GetLatest(ctx, req.Msg.GetNamespaceId(), req.Msg.GetSchemaId())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&stencilv1beta1.GetLatestSchemaResponse{
		Data: data,
	}), nil
}

func (a *API) GetSchema(ctx context.Context, req *connect.Request[stencilv1beta1.GetSchemaRequest]) (*connect.Response[stencilv1beta1.GetSchemaResponse], error) {
	_, data, err := a.schema.Get(ctx, req.Msg.GetNamespaceId(), req.Msg.GetSchemaId(), req.Msg.GetVersionId())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&stencilv1beta1.GetSchemaResponse{
		Data: data,
	}), nil
}

func (a *API) ListVersions(ctx context.Context, req *connect.Request[stencilv1beta1.ListVersionsRequest]) (*connect.Response[stencilv1beta1.ListVersionsResponse], error) {
	versions, err := a.schema.ListVersions(ctx, req.Msg.GetNamespaceId(), req.Msg.GetSchemaId())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&stencilv1beta1.ListVersionsResponse{Versions: versions}), nil
}

func (a *API) GetSchemaMetadata(ctx context.Context, req *connect.Request[stencilv1beta1.GetSchemaMetadataRequest]) (*connect.Response[stencilv1beta1.GetSchemaMetadataResponse], error) {
	meta, err := a.schema.GetMetadata(ctx, req.Msg.GetNamespaceId(), req.Msg.GetSchemaId())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&stencilv1beta1.GetSchemaMetadataResponse{
		Format:        stencilv1beta1.Schema_Format(stencilv1beta1.Schema_Format_value[meta.Format]),
		Compatibility: stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[meta.Compatibility]),
		Authority:     meta.Authority,
	}), nil
}

func (a *API) UpdateSchemaMetadata(ctx context.Context, req *connect.Request[stencilv1beta1.UpdateSchemaMetadataRequest]) (*connect.Response[stencilv1beta1.UpdateSchemaMetadataResponse], error) {
	meta, err := a.schema.UpdateMetadata(ctx, req.Msg.GetNamespaceId(), req.Msg.GetSchemaId(), &schema.Metadata{
		Compatibility: req.Msg.GetCompatibility().String(),
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&stencilv1beta1.UpdateSchemaMetadataResponse{
		Format:        stencilv1beta1.Schema_Format(stencilv1beta1.Schema_Format_value[meta.Format]),
		Compatibility: stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[meta.Compatibility]),
		Authority:     meta.Authority,
	}), nil
}

func (a *API) DeleteSchema(ctx context.Context, req *connect.Request[stencilv1beta1.DeleteSchemaRequest]) (*connect.Response[stencilv1beta1.DeleteSchemaResponse], error) {
	err := a.schema.Delete(ctx, req.Msg.GetNamespaceId(), req.Msg.GetSchemaId())
	message := "success"
	if err != nil {
		message = "failed"
	}
	return connect.NewResponse(&stencilv1beta1.DeleteSchemaResponse{
		Message: message,
	}), err
}

func (a *API) DeleteVersion(ctx context.Context, req *connect.Request[stencilv1beta1.DeleteVersionRequest]) (*connect.Response[stencilv1beta1.DeleteVersionResponse], error) {
	err := a.schema.DeleteVersion(ctx, req.Msg.GetNamespaceId(), req.Msg.GetSchemaId(), req.Msg.GetVersionId())
	message := "success"
	if err != nil {
		message = "failed"
	}
	return connect.NewResponse(&stencilv1beta1.DeleteVersionResponse{
		Message: message,
	}), err
}
