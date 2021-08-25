package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/odpf/stencil/server/api"
	"github.com/odpf/stencil/server/api/v1/genproto"
	"github.com/odpf/stencil/server/config"
	"github.com/odpf/stencil/server/proto"
	"github.com/odpf/stencil/server/snapshot"
	"github.com/odpf/stencil/server/store"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

// Router returns server router
func Router(api *api.API, config *config.Config) *gin.Engine {
	router := gin.New()
	addMiddleware(router, config)
	registerCustomValidations(router)
	registerRoutes(router, api)
	return router
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
	router := Router(api, config)

	// Create a gRPC server object
	s := grpc.NewServer()
	genproto.RegisterSnapshotServiceServer(s, api)
	genproto.RegisterStencilServiceServer(s, api)
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("0.0.0.0%s", port),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	genproto.RegisterSnapshotServiceHandler(ctx, gwmux, conn)
	gwmux.HandlePath("GET", "/ping", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		router.ServeHTTP(w, r)
	})
	gwmux.HandlePath("GET", "/v1/namespaces/{namespace}/descriptors/{name}/versions/{version}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		router.ServeHTTP(w, r)
	})
	gwmux.HandlePath("POST", "/v1/namespaces/{namespace}/descriptors", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		router.ServeHTTP(w, r)
	})
	// err = gen.RegisterStencilHandler(context.Background(), gwmux, conn)
	basemux := http.NewServeMux()
	basemux.Handle("/", gwmux)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	runWithGracefulShutdown(config, grpcHandlerFunc(s, basemux), func() {
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
