package router

import (
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/api/handlers"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/api/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Setup creates and configures the HTTP router.
func Setup(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Base middlewares
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Metrics)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// System endpoints
	r.Get("/healthz", h.Healthz)
	r.Get("/readyz", h.Readyz)
	r.Handle("/metrics", promhttp.Handler())

	// Swagger docs
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
	))

	// API endpoints
	r.Post("/analyze", h.Analyze)

	return r
}
