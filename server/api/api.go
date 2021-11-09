package api

import (
	"context"

	"github.com/odpf/stencil/server/namespace"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"github.com/odpf/stencil/server/schema"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// NamespaceService Service interface for namespace
type NamespaceService interface {
	Create(context.Context, namespace.Namespace) (namespace.Namespace, error)
	Update(context.Context, namespace.Namespace) (namespace.Namespace, error)
	List(context.Context) ([]string, error)
	Get(context.Context, string) (namespace.Namespace, error)
	Delete(context.Context, string) error
}

//SchemaService Service interface for schema management
type SchemaService interface {
	Create(context.Context, string, string, *schema.Metadata, []byte) (schema.SchemaInfo, error)
	List(context.Context, string) ([]string, error)
	Get(context.Context, string, string, int32) ([]byte, error)
	GetLatest(context.Context, string, string) ([]byte, error)
	ListVersions(context.Context, string, string) ([]int32, error)
	GetMetadata(context.Context, string, string) (*schema.Metadata, error)
	UpdateMetadata(context.Context, string, string, *schema.Metadata) (*schema.Metadata, error)
	DeleteVersion(context.Context, string, string, int32) error
}

//API holds all handlers
type API struct {
	stencilv1.UnimplementedStencilServiceServer
	grpc_health_v1.UnimplementedHealthServer
	Namespace NamespaceService
	Schema    SchemaService
}
