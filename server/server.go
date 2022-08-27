package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgraph-io/ristretto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	"github.com/odpf/stencil/config"
	"github.com/odpf/stencil/internal/store/postgres"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/odpf/stencil/core/namespace"
	"github.com/odpf/stencil/core/schema"
	"github.com/odpf/stencil/core/schema/provider"
	"github.com/odpf/stencil/core/search"
	"github.com/odpf/stencil/server/api"
	"github.com/odpf/stencil/server/logger"
	stencilv1beta1 "github.com/odpf/stencil/server/odpf/stencil/v1beta1"
	"github.com/odpf/stencil/server/validator"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Start Entry point to start the server
func Start(cfg config.Config) {
	ctx := context.Background()

	store := postgres.NewStore(cfg.DB.ConnectionString)
	namespaceService := namespace.Service{
		Repo: store,
	}
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,
		MaxCost:     cfg.CacheSizeInMB << 20,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	schemaService := schema.NewService(store, provider.NewSchemaProvider(), namespaceService, cache)
	searchService := &search.Service{
		Repo: store,
	}
	api := api.NewAPI(namespaceService, schemaService, searchService)

	port := fmt.Sprintf(":%s", cfg.Port)
	nr := getNewRelic(&cfg)
	mux := runtime.NewServeMux()

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
	api.RegisterSchemaHandlers(mux, nr)

	if err = stencilv1beta1.RegisterStencilServiceHandler(ctx, mux, conn); err != nil {
		log.Fatalln("Failed to register stencil service handler:", err)
	}

	runWithGracefulShutdown(&cfg, grpcHandlerFunc(s, mux), func() {
		conn.Close()
		s.GracefulStop()
		store.Close()
	})
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
