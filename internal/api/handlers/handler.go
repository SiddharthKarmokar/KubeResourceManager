package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/siddk/kube-resource-manager/internal/api/dto"
	domainmodels "github.com/siddk/kube-resource-manager/internal/domain/models"
	"github.com/siddk/kube-resource-manager/internal/errors"
	"github.com/siddk/kube-resource-manager/internal/services"
)

// Handler handles HTTP API requests.
type Handler struct {
	optimizer *services.OptimizerService
}

// NewHandler creates a new HTTP Handler.
func NewHandler(optimizer *services.OptimizerService) *Handler {
	return &Handler{
		optimizer: optimizer,
	}
}

// Healthz responds to health checks.
// @Summary Health check
// @Description Returns 200 if the server is alive
// @Tags System
// @Success 200 {string} string "OK"
// @Router /healthz [get]
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// Readyz responds to readiness checks.
// @Summary Readiness check
// @Description Returns 200 if the server is ready to accept traffic
// @Tags System
// @Success 200 {string} string "OK"
// @Router /readyz [get]
func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
	// In a real application, you might check dependencies here.
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// Analyze handles workload metrics analysis.
// @Summary Analyze workload metrics
// @Description Accepts a JSON array of workload metrics (CPU/memory request and average usage) and returns a JSON array of recommendations including confidence and severity.
// @Tags Optimization
// @Accept json
// @Produce json
// @Param request body []domainmodels.WorkloadMetrics true "Metrics batch (same schema as CLI JSON input)"
// @Success 200 {array} domainmodels.Recommendation
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /analyze [post]
func (h *Handler) Analyze(w http.ResponseWriter, r *http.Request) {
	var req dto.AnalyzeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, errors.NewAPIError(errors.CodeInvalidInput, "invalid json format"))
		return
	}

	recommendations, err := h.optimizer.AnalyzeBatch(r.Context(), req)
	if err != nil {
		if apiErr, ok := err.(errors.APIError); ok {
			h.writeError(w, apiErr)
			return
		}
		h.writeError(w, errors.NewAPIError(errors.CodeInternalError, "internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(recommendations)
}

func (h *Handler) writeError(w http.ResponseWriter, err errors.APIError) {
	w.Header().Set("Content-Type", "application/json")

	status := http.StatusInternalServerError
	if err.Code == errors.CodeInvalidInput || err.Code == errors.CodeInvalidRequest {
		status = http.StatusBadRequest
	}

	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errors.ErrorResponse{Error: err})
}

// Swag resolves request/response types referenced from handler comments; this keeps the import from being removed as unused.
var _ = domainmodels.WorkloadMetrics{}
