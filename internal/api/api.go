package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/odpf/stencil/core/search"
	"github.com/odpf/stencil/domain"
	stencilv1beta1 "github.com/odpf/stencil/proto/odpf/stencil/v1beta1"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type getSchemaData func(http.ResponseWriter, *http.Request, map[string]string) (*domain.Metadata, []byte, error)
type errHandleFunc func(http.ResponseWriter, *http.Request, map[string]string) error

// API holds all handlers
type API struct {
	stencilv1beta1.UnimplementedStencilServiceServer
	grpc_health_v1.UnimplementedHealthServer
	namespace domain.NamespaceService
	schema    domain.SchemaService
	search    search.SearchService
}

func NewAPI(namespace domain.NamespaceService, schema domain.SchemaService, search search.SearchService) *API {
	return &API{
		namespace: namespace,
		schema:    schema,
		search:    search,
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
