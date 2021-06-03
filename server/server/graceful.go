package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server/config"
)

func runWithGracefulShutdown(config *config.Config, router *gin.Engine, cleanUp func()) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 2)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		cleanUp()
		log.Fatal("Server forced to shutdown:", err)
	}
	cleanUp()

	log.Println("Server exiting")
}
