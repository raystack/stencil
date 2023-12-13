package api

import (
	"context"
	"fmt"
	newrelic2 "github.com/goto/stencil/pkg/newrelic"
	"net/http"
	"strconv"

	"github.com/goto/stencil/core/namespace"
	"github.com/goto/stencil/core/schema"
	"github.com/goto/stencil/core/search"
	stencilv1beta1 "github.com/goto/stencil/proto/v1beta1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/newrelic/go-agent/v3/newrelic"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type getSchemaData func(http.ResponseWriter, *http.Request, map[string]string) (*schema.Metadata, []byte, error)
type errHandleFunc func(http.ResponseWriter, *http.Request, map[string]string) error

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
	stencilv1beta1.UnimplementedStencilServiceServer
	grpc_health_v1.UnimplementedHealthServer
	namespace NamespaceService
	schema    SchemaService
	search    SearchService
	newrelic  newrelic2.Service
}

func NewAPI(namespace NamespaceService, schema SchemaService, search SearchService, nr newrelic2.Service) *API {
	return &API{
		namespace: namespace,
		schema:    schema,
		search:    search,
		newrelic:  nr,
	}
}

// RegisterSchemaHandlers registers HTTP handlers for schema download
func (a *API) RegisterSchemaHandlers(mux *runtime.ServeMux, app *newrelic.Application) {
	mux.HandlePath("GET", "/ping", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		fmt.Fprint(w, "pong")
	})
	mux.HandlePath(wrapHandler(app, "GET", "/v1beta1/namespaces/{namespace}/schemas/{name}/versions/{version}", handleSchemaResponse(mux, a.HTTPGetSchema)))
	mux.HandlePath(wrapHandler(app, "GET", "/v1beta1/namespaces/{namespace}/schemas/{name}", handleSchemaResponse(mux, a.HTTPLatestSchema)))
	mux.HandlePath(wrapHandler(app, "POST", "/v1beta1/namespaces/{namespace}/schemas/{name}", wrapErrHandler(mux, a.HTTPUpload)))
	mux.HandlePath(wrapHandler(app, "POST", "/v1beta1/namespaces/{namespace}/schemas/{name}/check", wrapErrHandler(mux, a.HTTPCheckCompatibility)))
}

func handleSchemaResponse(mux *runtime.ServeMux, getSchemaFn getSchemaData) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		meta, data, err := getSchemaFn(w, r, pathParams)
		if err != nil {
			_, outbound := runtime.MarshalerForRequest(mux, r)
			runtime.HTTPError(r.Context(), mux, outbound, w, r, err)
			return
		}
		contentType := "application/json"
		if meta.Format == "FORMAT_PROTOBUF" {
			contentType = "application/octet-stream"
		}
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func wrapErrHandler(mux *runtime.ServeMux, handler errHandleFunc) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		err := handler(w, r, pathParams)
		if err != nil {
			_, outbound := runtime.MarshalerForRequest(mux, r)
			runtime.DefaultHTTPErrorHandler(r.Context(), mux, outbound, w, r, err)
			return
		}
	}
}

func wrapHandler(app *newrelic.Application, method, pattern string, handler runtime.HandlerFunc) (string, string, runtime.HandlerFunc) {
	if app == nil {
		return method, pattern, handler
	}
	return method, pattern, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		txn := app.StartTransaction(method + " " + pattern)
		defer txn.End()
		w = txn.SetWebResponse(w)
		txn.SetWebRequestHTTP(r)
		r = newrelic.RequestWithTransactionContext(r, txn)
		handler(w, r, pathParams)
	}
}
