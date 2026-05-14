package handlers

// prometheusMetricsOpenAPI documents GET /metrics for Swagger generation. The live handler is promhttp.Handler() wired in the router.
//
// @Summary Prometheus metrics
// @Description Prometheus text exposition format (OpenMetrics / text/plain).
// @Tags System
// @Produce plain
// @Success 200 {string} string "Prometheus metrics"
// @Router /metrics [get]
func prometheusMetricsOpenAPI() {}

var _ = prometheusMetricsOpenAPI
