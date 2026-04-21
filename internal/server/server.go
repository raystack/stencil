package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	connectcors "connectrpc.com/cors"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/validate"
	"github.com/dgraph-io/ristretto"
	"github.com/raystack/salt/server/spa"
	"github.com/raystack/stencil/config"
	"github.com/raystack/stencil/core/namespace"
	"github.com/raystack/stencil/core/schema"
	"github.com/raystack/stencil/core/schema/provider"
	"github.com/raystack/stencil/core/search"
	stencilv1beta1connect "github.com/raystack/stencil/gen/raystack/stencil/v1beta1/stencilv1beta1connect"
	"github.com/raystack/stencil/internal/api"
	"github.com/raystack/stencil/internal/middleware"
	"github.com/raystack/stencil/internal/store/postgres"
	"github.com/raystack/stencil/ui"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// Start Entry point to start the server
func Start(cfg config.Config) {
	logger := slog.Default().With("component", "server")

	db := postgres.NewStore(cfg.DB.ConnectionString)
	defer db.Close()

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

	handler := api.NewAPI(namespaceService, schemaService, searchService)

	// Build interceptor chain
	validateInterceptor := validate.NewInterceptor()
	interceptors := connect.WithInterceptors(
		middleware.Recovery(),
		middleware.Logger(),
		validateInterceptor,
		middleware.ErrorResponse(),
	)

	// Create HTTP mux
	mux := http.NewServeMux()

	// Register connect service handler
	path, svcHandler := stencilv1beta1connect.NewStencilServiceHandler(
		handler,
		interceptors,
		connect.WithReadMaxBytes(cfg.MaxRecvMsgSize),
		connect.WithSendMaxBytes(cfg.MaxSendMsgSize),
	)
	mux.Handle(path, svcHandler)

	// Register gRPC reflection for tooling compatibility (grpcurl, etc.)
	reflector := grpcreflect.NewStaticReflector(
		"raystack.stencil.v1beta1.StencilService",
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// Register custom HTTP handlers for schema binary endpoints
	handler.RegisterSchemaHandlers(mux)

	// Health check endpoint
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})

	// UI SPA handler
	spaHandler, err := spa.Handler(ui.Assets, "build", "index.html", false)
	if err != nil {
		log.Fatalln("Failed to load spa:", err)
	}
	mux.Handle("/ui/", http.StripPrefix("/ui", spaHandler))

	// CORS middleware
	allowedOrigins := cfg.CORS.AllowedOrigins
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	}
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})

	// Create HTTP server with h2c support for HTTP/2 without TLS
	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      h2c.NewHandler(corsHandler.Handler(mux), &http2.Server{}),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("starting server", "addr", addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	logger.Info("server stopped")
}
