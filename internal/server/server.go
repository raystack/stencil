package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgraph-io/ristretto"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	"github.com/odpf/salt/spa"
	"github.com/odpf/stencil/config"
	"github.com/odpf/stencil/internal/store/postgres"
	"github.com/odpf/stencil/ui"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/odpf/stencil/core/namespace"
	"github.com/odpf/stencil/core/schema"
	"github.com/odpf/stencil/core/schema/provider"
	"github.com/odpf/stencil/core/search"
	"github.com/odpf/stencil/internal/api"
	"github.com/odpf/stencil/pkg/logger"
	"github.com/odpf/stencil/pkg/validator"
	stencilv1beta1 "github.com/odpf/stencil/proto/odpf/stencil/v1beta1"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Start Entry point to start the server
func Start(cfg config.Config) {
	ctx := context.Background()

	db := postgres.NewStore(cfg.DB.ConnectionString)

	namespaceRepository := postgres.NewNamespaceRepository(db)
	namespaceService := namespace.NewService(namespaceRepository)

	schemaRepository := postgres.NewSchemaRepository(db)
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,
		MaxCost:     cfg.CacheSizeInMB << 20,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	schemaService := schema.NewService(schemaRepository, provider.NewSchemaProvider(), namespaceService, cache)

	searchRepository := postgres.NewSearchRepository(db)
	searchService := search.NewService(searchRepository)

	api := api.NewAPI(namespaceService, schemaService, searchService)

	port := fmt.Sprintf(":%s", cfg.Port)
	nr := getNewRelic(&cfg)
	gatewayMux := runtime.NewServeMux()

	// init grpc server
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			nrgrpc.UnaryServerInterceptor(nr),
			grpc_zap.UnaryServerInterceptor(logger.Logger),
			validator.UnaryServerInterceptor())),
		grpc.MaxRecvMsgSize(cfg.GRPC.MaxRecvMsgSizeInMB << 20),
		grpc.MaxSendMsgSize(cfg.GRPC.MaxSendMsgSizeInMB << 20),
	}
	// Create a gRPC server object
	s := grpc.NewServer(opts...)
	stencilv1beta1.RegisterStencilServiceServer(s, api)
	grpc_health_v1.RegisterHealthServer(s, api)
	conn, err := grpc.DialContext(
		context.Background(),
		port,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	api.RegisterSchemaHandlers(gatewayMux, nr)

	if err = stencilv1beta1.RegisterStencilServiceHandler(ctx, gatewayMux, conn); err != nil {
		log.Fatalln("Failed to register stencil service handler:", err)
	}

	rtr := mux.NewRouter()

	spaHandler, err := spa.Handler(ui.Assets, "build", "index.html", false)
	if err != nil {
		log.Fatalln("Failed to load spa:", err)
	}
	rtr.PathPrefix("/ui").Handler(http.StripPrefix("/ui", spaHandler))

	runWithGracefulShutdown(&cfg, grpcHandlerFunc(s, gatewayMux, rtr), func() {
		conn.Close()
		s.GracefulStop()
		db.Close()
	})
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler, uiHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			if strings.HasPrefix(r.URL.Path, "/ui") {
				uiHandler.ServeHTTP(w, r)
			} else {
				otherHandler.ServeHTTP(w, r)
			}
		}
	}), &http2.Server{})
}
