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

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/odpf/stencil/server/api"
	"github.com/odpf/stencil/server/api/v1/pb"
	"github.com/odpf/stencil/server/config"
	"github.com/odpf/stencil/server/logger"
	"github.com/odpf/stencil/server/proto"
	"github.com/odpf/stencil/server/snapshot"
	"github.com/odpf/stencil/server/store"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
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
func Start() {
	ctx := context.Background()
	config := config.LoadConfig()
	db := store.NewDBStore(config)
	stRepo := snapshot.NewSnapshotRepository(db)
	protoRepo := proto.NewProtoRepository(db)
	protoService := proto.NewService(protoRepo, stRepo)
	api := &api.API{
		Store:    protoService,
		Metadata: stRepo,
	}
	port := fmt.Sprintf(":%s", config.Port)
	nr := getNewRelic(config)
	mux := Router(api, config)

	// init grpc server
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			nrgrpc.UnaryServerInterceptor(nr),
			grpc_zap.UnaryServerInterceptor(logger.Logger))),
		grpc.MaxRecvMsgSize(config.GRPCMaxRecvMsgSizeInMB << 20),
		grpc.MaxSendMsgSize(config.GRPCMaxSendMsgSizeInMB << 20),
	}
	// Create a gRPC server object
	s := grpc.NewServer(opts...)
	pb.RegisterStencilServiceServer(s, api)
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("0.0.0.0%s", port),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	pb.RegisterStencilServiceHandler(ctx, mux, conn)

	runWithGracefulShutdown(config, grpcHandlerFunc(s, mux), func() {
		conn.Close()
		s.GracefulStop()
		db.Close()
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
