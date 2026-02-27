package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/didikurnia/api-quest/internal/config"
	"github.com/didikurnia/api-quest/internal/router"
	"github.com/didikurnia/api-quest/internal/store"
)

func main() {
	cfg := config.Load()
	bookStore := store.NewBookStore()

	r := router.Setup(cfg, bookStore)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 API Quest server running on http://localhost:%s\n", cfg.Port)
		log.Printf("📚 API Docs available at http://localhost:%s/docs\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited gracefully")
}
