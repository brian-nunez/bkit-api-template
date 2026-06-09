package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/brian-nunez/bhttp/pkg/bsuite"
	"github.com/brian-nunez/bkit-api-template/views/pages"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
)

// TemplateRenderer implements the echo.Renderer interface for Templ templates
type TemplateRenderer struct{}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	component, ok := data.(templ.Component)
	if !ok {
		return fmt.Errorf("data is not a templ.Component")
	}
	return component.Render(c.Request().Context(), w)
}

type Server struct {
	container *bsuite.Service
	startTime time.Time
}

func New(container *bsuite.Service) *Server {
	return &Server{
		container: container,
		startTime: time.Now(),
	}
}

func (s *Server) Run(ctx context.Context) error {
	cfg := s.container.Config()
	port := cfg.Int("server.port")
	if port == 0 {
		port = 8080
	}

	e := echo.New()
	e.Renderer = &TemplateRenderer{}

	// Global Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	// OpenTelemetry Tracing Middleware if enabled
	if cfg.Bool("telemetry.enabled") && cfg.Bool("telemetry.enable_trace") {
		serviceName := cfg.String("telemetry.service_name")
		e.Use(otelecho.Middleware(
			serviceName,
			otelecho.WithPropagators(otel.GetTextMapPropagator()),
			otelecho.WithTracerProvider(otel.GetTracerProvider()),
		))
	}

	// Serve Static Assets
	e.Static("/assets", "./assets")

	// Routes
	e.GET("/", s.handleHome)
	e.GET("/api/status", s.handleStatusGrid)

	// Metrics endpoint (Prometheus pull)
	if s.container.Telemetry() != nil && s.container.Telemetry().HTTPHandler() != nil {
		e.GET("/metrics", echo.WrapHandler(s.container.Telemetry().HTTPHandler()))
	}

	// Start server inside goroutine
	errChan := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf(":%d", port)
		e.Logger.Infof("HTTP server listening on %s", addr)
		if err := e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	// Graceful shutdown wait
	select {
	case <-ctx.Done():
		e.Logger.Info("Shutting down HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return e.Shutdown(shutdownCtx)
	case err := <-errChan:
		return err
	}
}

// collectStatus queries DB and KV status for the dashboard
func (s *Server) collectStatus(ctx context.Context) pages.ServiceStatus {
	cfg := s.container.Config()

	status := pages.ServiceStatus{
		ServiceName:        cfg.String("telemetry.service_name"),
		Environment:        cfg.String("telemetry.environment"),
		DBDriver:           cfg.String("db.driver"),
		KVDriver:           cfg.String("kv.driver"),
		Uptime:             time.Since(s.startTime).Round(time.Second).String(),
		TelemetryConnected: s.container.Telemetry() != nil,
	}

	// Check DB
	if s.container.DB() != nil {
		dbCtx, dbCancel := context.WithTimeout(ctx, 2*time.Second)
		defer dbCancel()
		if err := s.container.DB().Ping(dbCtx); err != nil {
			status.DBConnected = false
			status.DBError = err.Error()
		} else {
			status.DBConnected = true
		}
	} else {
		status.DBConnected = false
		status.DBError = "Database service disabled in config"
	}

	// Check KV
	if s.container.KV() != nil {
		kvCtx, kvCancel := context.WithTimeout(ctx, 2*time.Second)
		defer kvCancel()
		if err := s.container.KV().HealthCheck(kvCtx); err != nil {
			status.KVConnected = false
			status.KVError = err.Error()
		} else {
			status.KVConnected = true
		}
	} else {
		status.KVConnected = false
		status.KVError = "KV store service disabled in config"
	}

	return status
}

func (s *Server) handleHome(c echo.Context) error {
	status := s.collectStatus(c.Request().Context())
	return c.Render(http.StatusOK, "", pages.Home(status))
}

func (s *Server) handleStatusGrid(c echo.Context) error {
	status := s.collectStatus(c.Request().Context())
	return c.Render(http.StatusOK, "", pages.StatusGrid(status))
}
