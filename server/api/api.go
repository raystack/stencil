package api

import (
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/odpf/stencil/server/domain"
	stencilv1beta1 "github.com/odpf/stencil/server/odpf/stencil/v1beta1"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type getSchemaData func(http.ResponseWriter, *http.Request, map[string]string) (*domain.Metadata, []byte, error)
type errHandleFunc func(http.ResponseWriter, *http.Request, map[string]string) error

//API holds all handlers
type API struct {
	stencilv1beta1.UnimplementedStencilServiceServer
	grpc_health_v1.UnimplementedHealthServer
	Namespace     domain.NamespaceService
	Schema        domain.SchemaService
	SearchService domain.SearchService
}

// RegisterSchemaHandlers registers HTTP handlers for schema download
func (a *API) RegisterSchemaHandlers(mux *runtime.ServeMux) {
	mux.HandlePath("GET", "/ping", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		fmt.Fprint(w, "pong")
	})
	mux.HandlePath("GET", "/v1beta1/namespaces/{namespace}/schemas/{name}/versions/{version}", handleSchemaResponse(mux, a.HTTPGetSchema))
	mux.HandlePath("GET", "/v1beta1/namespaces/{namespace}/schemas/{name}", handleSchemaResponse(mux, a.HTTPLatestSchema))
	mux.HandlePath("POST", "/v1beta1/namespaces/{namespace}/schemas/{name}", wrapHandler(mux, a.HTTPUpload))
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
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func wrapHandler(mux *runtime.ServeMux, handler errHandleFunc) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		err := handler(w, r, pathParams)
		if err != nil {
			_, outbound := runtime.MarshalerForRequest(mux, r)
			runtime.DefaultHTTPErrorHandler(r.Context(), mux, outbound, w, r, err)
			return
		}
	}
}
