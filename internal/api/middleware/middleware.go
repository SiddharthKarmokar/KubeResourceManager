package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/logger"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/observability"
)

// Metrics records Prometheus metrics for each request.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()
		status := http.StatusText(ww.Status())

		observability.RequestCount.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		observability.RequestLatency.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}

// Logger injects a request-scoped logger and logs request details.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())
		if reqID == "" {
			reqID = "unknown"
		}

		log := logger.FromContext(r.Context()).With(
			"request_id", reqID,
			"method", r.Method,
			"path", r.URL.Path,
		)

		ctx := logger.WithLogger(r.Context(), log)

		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r.WithContext(ctx))

		log.Info("request completed",
			"status", ww.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

// Timeout adds a timeout to the request context.
func Timeout(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
