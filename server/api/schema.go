package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/odpf/stencil/server/domain"
	stencilv1beta1 "github.com/odpf/stencil/server/odpf/stencil/v1beta1"
)

func (a *API) CreateSchema(ctx context.Context, in *stencilv1beta1.CreateSchemaRequest) (*stencilv1beta1.CreateSchemaResponse, error) {
	metadata := &domain.Metadata{Format: in.GetFormat().String(), Compatibility: in.GetCompatibility().String()}
	sc, err := a.Schema.Create(ctx, in.NamespaceId, in.SchemaId, metadata, in.GetData())
	return &stencilv1beta1.CreateSchemaResponse{
		Version:  sc.Version,
		Id:       sc.ID,
		Location: sc.Location,
	}, err
}
func (a *API) HTTPUpload(w http.ResponseWriter, req *http.Request, pathParams map[string]string) error {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	format := req.Header.Get("X-Format")
	compatibility := req.Header.Get("X-Compatibility")
	metadata := &domain.Metadata{Format: format, Compatibility: compatibility}
	namespaceID := pathParams["namespace"]
	schemaName := pathParams["name"]
	sc, err := a.Schema.Create(req.Context(), namespaceID, schemaName, metadata, data)
	if err != nil {
		return err
	}
	respData, _ := json.Marshal(sc)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(respData)
	return nil
}

func (a *API) ListSchemas(ctx context.Context, in *stencilv1beta1.ListSchemasRequest) (*stencilv1beta1.ListSchemasResponse, error) {
	schemas, err := a.Schema.List(ctx, in.Id)
	return &stencilv1beta1.ListSchemasResponse{Schemas: schemas}, err
}

func (a *API) GetLatestSchema(ctx context.Context, in *stencilv1beta1.GetLatestSchemaRequest) (*stencilv1beta1.GetLatestSchemaResponse, error) {
	_, data, err := a.Schema.GetLatest(ctx, in.NamespaceId, in.SchemaId)
	return &stencilv1beta1.GetLatestSchemaResponse{
		Data: data,
	}, err
}

func (a *API) HTTPLatestSchema(w http.ResponseWriter, req *http.Request, pathParams map[string]string) (*domain.Metadata, []byte, error) {
	namespaceID := pathParams["namespace"]
	schemaName := pathParams["name"]
	return a.Schema.GetLatest(req.Context(), namespaceID, schemaName)
}

func (a *API) GetSchema(ctx context.Context, in *stencilv1beta1.GetSchemaRequest) (*stencilv1beta1.GetSchemaResponse, error) {
	_, data, err := a.Schema.Get(ctx, in.NamespaceId, in.SchemaId, in.GetVersionId())
	return &stencilv1beta1.GetSchemaResponse{
		Data: data,
	}, err
}

func (a *API) HTTPGetSchema(w http.ResponseWriter, req *http.Request, pathParams map[string]string) (*domain.Metadata, []byte, error) {
	namespaceID := pathParams["namespace"]
	schemaName := pathParams["name"]
	versionString := pathParams["version"]
	v, err := strconv.ParseInt(versionString, 10, 32)
	if err != nil {
		return nil, nil, &runtime.HTTPStatusError{HTTPStatus: http.StatusBadRequest, Err: errors.New("invalid version number")}
	}
	return a.Schema.Get(req.Context(), namespaceID, schemaName, int32(v))
}

func (a *API) ListVersions(ctx context.Context, in *stencilv1beta1.ListVersionsRequest) (*stencilv1beta1.ListVersionsResponse, error) {
	versions, err := a.Schema.ListVersions(ctx, in.NamespaceId, in.SchemaId)
	return &stencilv1beta1.ListVersionsResponse{Versions: versions}, err
}

func (a *API) GetSchemaMetadata(ctx context.Context, in *stencilv1beta1.GetSchemaMetadataRequest) (*stencilv1beta1.GetSchemaMetadataResponse, error) {
	meta, err := a.Schema.GetMetadata(ctx, in.NamespaceId, in.SchemaId)
	return &stencilv1beta1.GetSchemaMetadataResponse{
		Format:        stencilv1beta1.Schema_Format(stencilv1beta1.Schema_Format_value[meta.Format]),
		Compatibility: stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[meta.Compatibility]),
		Authority:     meta.Authority,
	}, err
}

func (a *API) UpdateSchemaMetadata(ctx context.Context, in *stencilv1beta1.UpdateSchemaMetadataRequest) (*stencilv1beta1.UpdateSchemaMetadataResponse, error) {
	meta, err := a.Schema.UpdateMetadata(ctx, in.NamespaceId, in.SchemaId, &domain.Metadata{
		Compatibility: in.Compatibility.String(),
	})
	return &stencilv1beta1.UpdateSchemaMetadataResponse{
		Format:        stencilv1beta1.Schema_Format(stencilv1beta1.Schema_Format_value[meta.Format]),
		Compatibility: stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[meta.Compatibility]),
		Authority:     meta.Authority,
	}, err
}

func (a *API) DeleteSchema(ctx context.Context, in *stencilv1beta1.DeleteSchemaRequest) (*stencilv1beta1.DeleteSchemaResponse, error) {
	err := a.Schema.Delete(ctx, in.NamespaceId, in.SchemaId)
	message := "success"
	if err != nil {
		message = "failed"
	}
	return &stencilv1beta1.DeleteSchemaResponse{
		Message: message,
	}, err
}

func (a *API) DeleteVersion(ctx context.Context, in *stencilv1beta1.DeleteVersionRequest) (*stencilv1beta1.DeleteVersionResponse, error) {
	err := a.Schema.DeleteVersion(ctx, in.NamespaceId, in.SchemaId, in.GetVersionId())
	message := "success"
	if err != nil {
		message = "failed"
	}
	return &stencilv1beta1.DeleteVersionResponse{
		Message: message,
	}, err
}
