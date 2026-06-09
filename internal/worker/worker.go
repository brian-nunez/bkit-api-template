package worker

import (
	"context"
	"log"
	"time"

	"github.com/brian-nunez/bhttp/pkg/bsuite"
)

type Worker struct {
	container *bsuite.Service
}

func New(container *bsuite.Service) *Worker {
	return &Worker{
		container: container,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	log.Println("[Worker] Background task runner started. Running every 15s...")

	for {
		select {
		case <-ctx.Done():
			log.Println("[Worker] Shutting down background worker...")
			return nil
		case <-ticker.C:
			w.performCheck(ctx)
		}
	}
}

func (w *Worker) performCheck(ctx context.Context) {
	log.Println("[Worker] Performing periodic health checks & background sync...")

	// Verify DB health
	db := w.container.DB()
	if db != nil {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			log.Printf("[Worker] [DB WARNING] DB health check failed: %v", err)
		} else {
			log.Println("[Worker] [DB INFO] DB connection is healthy")
		}
	}

	// Verify KV health
	kv := w.container.KV()
	if kv != nil {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		if err := kv.HealthCheck(ctx); err != nil {
			log.Printf("[Worker] [KV WARNING] KV health check failed: %v", err)
		} else {
			log.Println("[Worker] [KV INFO] KV store connection is healthy")

			// Increment a counter in KV as an example of writes
			valStr, err := kv.Get(ctx, "worker:ticks")
			var ticks string
			if err != nil {
				ticks = "0"
			} else {
				ticks = valStr
			}
			log.Printf("[Worker] [KV INFO] Last recorded ticks: %s", ticks)

			// Update ticks (simulated)
			_ = kv.Set(ctx, "worker:ticks", time.Now().Format(time.RFC3339), 1*time.Hour)
		}
	}
}
