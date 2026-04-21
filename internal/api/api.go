package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/raystack/stencil/core/namespace"
	"github.com/raystack/stencil/core/schema"
	"github.com/raystack/stencil/core/search"
)

type NamespaceService interface {
	Create(ctx context.Context, ns namespace.Namespace) (namespace.Namespace, error)
	Update(ctx context.Context, ns namespace.Namespace) (namespace.Namespace, error)
	List(ctx context.Context) ([]namespace.Namespace, error)
	Get(ctx context.Context, name string) (namespace.Namespace, error)
	Delete(ctx context.Context, name string) error
}

type SchemaService interface {
	CheckCompatibility(ctx context.Context, nsName, schemaName, compatibility string, data []byte) error
	Create(ctx context.Context, nsName string, schemaName string, metadata *schema.Metadata, data []byte) (schema.SchemaInfo, error)
	Get(ctx context.Context, namespace string, schemaName string, version int32) (*schema.Metadata, []byte, error)
	Delete(ctx context.Context, namespace string, schemaName string) error
	DeleteVersion(ctx context.Context, namespace string, schemaName string, version int32) error
	GetLatest(ctx context.Context, namespace string, schemaName string) (*schema.Metadata, []byte, error)
	GetMetadata(ctx context.Context, namespace, schemaName string) (*schema.Metadata, error)
	UpdateMetadata(ctx context.Context, namespace, schemaName string, meta *schema.Metadata) (*schema.Metadata, error)
	List(ctx context.Context, namespaceID string) ([]schema.Schema, error)
	ListVersions(ctx context.Context, namespaceID string, schemaName string) ([]int32, error)
}

type SearchService interface {
	Search(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error)
}

type API struct {
	namespace NamespaceService
	schema    SchemaService
	search    SearchService
}

func NewAPI(namespace NamespaceService, schema SchemaService, search SearchService) *API {
	return &API{
		namespace: namespace,
		schema:    schema,
		search:    search,
	}
}

// RegisterSchemaHandlers registers custom HTTP handlers for schema binary endpoints.
// These serve raw schema bytes (protobuf/JSON) rather than standard RPC responses.
func (a *API) RegisterSchemaHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1beta1/namespaces/{namespace}/schemas/{name}/versions/{version}", a.handleGetSchema)
	mux.HandleFunc("GET /v1beta1/namespaces/{namespace}/schemas/{name}", a.handleGetLatestSchema)
	mux.HandleFunc("POST /v1beta1/namespaces/{namespace}/schemas/{name}", a.handleUploadSchema)
	mux.HandleFunc("POST /v1beta1/namespaces/{namespace}/schemas/{name}/check", a.handleCheckCompatibility)
}

func (a *API) handleGetSchema(w http.ResponseWriter, r *http.Request) {
	namespaceID := r.PathValue("namespace")
	schemaName := r.PathValue("name")
	versionStr := r.PathValue("version")

	v, err := strconv.ParseInt(versionStr, 10, 32)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid version number")
		return
	}

	meta, data, err := a.schema.Get(r.Context(), namespaceID, schemaName, int32(v))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeSchemaResponse(w, meta, data)
}

func (a *API) handleGetLatestSchema(w http.ResponseWriter, r *http.Request) {
	namespaceID := r.PathValue("namespace")
	schemaName := r.PathValue("name")

	meta, data, err := a.schema.GetLatest(r.Context(), namespaceID, schemaName)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeSchemaResponse(w, meta, data)
}

func (a *API) handleUploadSchema(w http.ResponseWriter, r *http.Request) {
	data, err := readBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	format := r.Header.Get("X-Format")
	compatibility := r.Header.Get("X-Compatibility")
	metadata := &schema.Metadata{Format: format, Compatibility: compatibility}
	namespaceID := r.PathValue("namespace")
	schemaName := r.PathValue("name")

	sc, err := a.schema.Create(r.Context(), namespaceID, schemaName, metadata, data)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sc)
}

func (a *API) handleCheckCompatibility(w http.ResponseWriter, r *http.Request) {
	data, err := readBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	compatibility := r.Header.Get("X-Compatibility")
	namespaceID := r.PathValue("namespace")
	schemaName := r.PathValue("name")

	if err := a.schema.CheckCompatibility(r.Context(), namespaceID, schemaName, compatibility, data); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func writeSchemaResponse(w http.ResponseWriter, meta *schema.Metadata, data []byte) {
	contentType := "application/json"
	if meta.Format == "FORMAT_PROTOBUF" {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func writeError(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeServiceError(w http.ResponseWriter, err error) {
	writeError(w, http.StatusInternalServerError, err.Error())
}

func readBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}
