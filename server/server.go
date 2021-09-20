package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	"github.com/odpf/stencil/config"
	"github.com/odpf/stencil/models"
	"github.com/odpf/stencil/search"
	"github.com/odpf/stencil/storage/postgres"
	"github.com/pkg/errors"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/odpf/stencil/server/api"
	"github.com/odpf/stencil/server/logger"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"github.com/odpf/stencil/server/proto"
	"github.com/odpf/stencil/server/snapshot"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	gproto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Router returns server router
func Router(api *api.API, config *config.Config) *runtime.ServeMux {
	gwmux := runtime.NewServeMux()
	router := gin.New()
	addMiddleware(router, config)
	registerCustomValidations(router)
	registerRoutes(router, gwmux, api)
	return gwmux
}

// Start Entry point to start the server
func Start(cfg config.Config) {
	ctx := context.Background()

	store := postgres.NewStore(cfg.DB.ConnectionString)
	protoService := proto.NewService(store)
	metaService := snapshot.NewService(store)
	cache := search.NewInMemoryStore()
	api := &api.API{
		Store:         protoService,
		Metadata:      metaService,
		SearchService: cache,
	}
	port := fmt.Sprintf(":%s", cfg.Port)
	nr := getNewRelic(&cfg)
	mux := Router(api, &cfg)

	// init grpc server
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			nrgrpc.UnaryServerInterceptor(nr),
			grpc_zap.UnaryServerInterceptor(logger.Logger))),
		grpc.MaxRecvMsgSize(cfg.GRPC.MaxRecvMsgSizeInMB << 20),
		grpc.MaxSendMsgSize(cfg.GRPC.MaxSendMsgSizeInMB << 20),
	}
	// Create a gRPC server object
	s := grpc.NewServer(opts...)
	stencilv1.RegisterStencilServiceServer(s, api)
	grpc_health_v1.RegisterHealthServer(s, api)
	conn, err := grpc.DialContext(
		context.Background(),
		port,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	stencilv1.RegisterStencilServiceHandler(ctx, mux, conn)

	err = buildSchemaIndex(ctx, cache, api)
	if err != nil {
		panic(err)
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

func buildSchemaIndex(ctx context.Context, cache *search.InMemoryStore, api *api.API) error {

	snapshots, err := api.Metadata.List(ctx, &models.Snapshot{})

	if err != nil {
		return errors.Wrap(err, "error getting snapshots")
	}
	for _, ss := range snapshots {

		descriptorSetBytes, err := api.Store.Get(ctx, ss, []string{})
		if err != nil {
			return errors.Wrap(err, "error getting descriptor set from store")
		}

		fds := &descriptorpb.FileDescriptorSet{}
		err = gproto.Unmarshal(descriptorSetBytes, fds)
		if err != nil {
			return errors.Wrap(err, "error unmarshalling descriptor set proto")
		}

		for _, proto := range fds.File {
			for _, m := range proto.GetMessageType() {
				fields := make([]string, 0)

				for _, f := range m.Field {
					fields = append(fields, f.GetName())
				}
				if err := cache.Index(ctx, &search.IndexRequest{
					Namespace: ss.Namespace,
					Version:   ss.Version,
					Fields:    fields,
					Name:      ss.Name,
					Latest:    ss.Latest,
					Message:   m.GetName(),
					Package:   proto.GetPackage(),
				}); err != nil {
					return errors.Wrap(err, "error indexing fields for search")
				}
			}

		}

	}

	return nil

}
