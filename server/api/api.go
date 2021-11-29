package api

import (
	"github.com/odpf/stencil/server/domain"
	stencilv1beta1 "github.com/odpf/stencil/server/odpf/stencil/v1beta1"
	"google.golang.org/grpc/health/grpc_health_v1"
)

//API holds all handlers
type API struct {
	stencilv1beta1.UnimplementedStencilServiceServer
	grpc_health_v1.UnimplementedHealthServer
	Namespace domain.NamespaceService
	Schema    domain.SchemaService
}
