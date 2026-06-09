package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brian-nunez/bhttp/pkg/brun"
	"github.com/brian-nunez/bhttp/pkg/bsuite"
	"github.com/brian-nunez/bkit-api-template/internal/config"
	"github.com/brian-nunez/bkit-api-template/internal/server"
	"github.com/brian-nunez/bkit-api-template/internal/worker"
)

func main() {
	// 1. Create root context with signal cancellation
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Initializing application...")

	// 2. Load configurations using bconfig
	cfg, err := config.Load(ctx)
	if err != nil {
		log.Fatalf("Fatal: failed to load configuration: %v", err)
	}

	// 3. Initialize BSuite service container (db, kv, telemetry)
	service, err := bsuite.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Fatal: failed to initialize bsuite container: %v", err)
	}
	defer func() {
		log.Println("Shutting down bsuite container...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := service.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error during bsuite container shutdown: %v", err)
		}
		log.Println("Container shutdown complete.")
	}()

	log.Printf("Services initialized. Service: %q, Env: %q",
		cfg.String("telemetry.service_name"), cfg.String("telemetry.environment"))

	// 4. Initialize API Server and Background Worker
	apiServer := server.New(service)
	bgWorker := worker.New(service)

	// 5. Initialize brun Manager and register runnables
	manager := brun.New()
	manager.Register(apiServer, bgWorker)

	log.Println("Starting concurrent runnables...")

	// 6. Start the manager (runs concurrently until context cancel or task failure)
	if err := manager.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("Fatal: manager exited with error: %v", err)
	}

	log.Println("Application stopped cleanly.")
}
