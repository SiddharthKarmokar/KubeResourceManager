package recommendation

import (
	"context"
	"fmt"
	"math"

	"github.com/siddk/kube-resource-manager/internal/config"
	"github.com/siddk/kube-resource-manager/internal/domain/models"
	"github.com/siddk/kube-resource-manager/internal/engine/heuristics"
	"github.com/siddk/kube-resource-manager/internal/engine/scoring"
)

// Recommender handles the generation of optimization recommendations.
type Recommender struct {
	cfg    config.RecommendationConfig
	scorer *scoring.Scorer
}

// NewRecommender creates a new Recommender instance.
func NewRecommender(cfg config.RecommendationConfig, scorer *scoring.Scorer) *Recommender {
	return &Recommender{
		cfg:    cfg,
		scorer: scorer,
	}
}

// GenerateRecommendation evaluates a single workload and produces a recommendation.
func (r *Recommender) GenerateRecommendation(ctx context.Context, metrics models.WorkloadMetrics) models.Recommendation {
	draftCPU := heuristics.CalculateSafeCPU(ctx, metrics.CPUUsageAvg, r.cfg)
	draftMem := heuristics.CalculateSafeMemory(ctx, metrics.MemoryUsageAvg, r.cfg)

	// Do not recommend increasing resources if the current request is already lower than the recommendation
	if draftCPU > metrics.CPURequest {
		draftCPU = metrics.CPURequest
	}
	if draftMem > metrics.MemoryRequest {
		draftMem = metrics.MemoryRequest
	}

	recCPU, recMem, degraded := finalizeRecommendations(draftCPU, draftMem, metrics, r.cfg)

	cpuReductionPct := 0
	if metrics.CPURequest > 0 {
		cpuReductionPct = int(math.Round(float64(metrics.CPURequest-recCPU) / float64(metrics.CPURequest) * 100))
		if cpuReductionPct < 0 {
			cpuReductionPct = 0
		}
	}

	memReductionPct := 0
	if metrics.MemoryRequest > 0 {
		memReductionPct = int(math.Round(float64(metrics.MemoryRequest-recMem) / float64(metrics.MemoryRequest) * 100))
		if memReductionPct < 0 {
			memReductionPct = 0
		}
	}

	confidence := r.scorer.CalculateConfidence(ctx, metrics, recCPU, recMem, degraded)
	severity := r.scorer.CalculateSeverity(ctx, cpuReductionPct, memReductionPct, degraded)

	reason := "Usage is close to requested resources; no change recommended."
	if cpuReductionPct > 0 || memReductionPct > 0 {
		reason = "Average usage significantly below requested resources with stable utilization patterns."
	}
	if degraded {
		reason = "Recommendation required safety clamps; verify metrics quality and configuration."
	}

	var warnings []string
	warnings = append(warnings, fmt.Sprintf("Recommendation includes %.1fx CPU and %.1fx Memory safety buffer.", r.cfg.CPUSafetyBuffer, r.cfg.MemorySafetyBuffer))
	warnings = append(warnings, "Current implementation assumes point-in-time average utilization data only.")
	if degraded {
		warnings = append(warnings, "Emergency invariants were applied (non-zero request, at or above average usage, minimum floors).")
	}

	savingsCPU := fmt.Sprintf("%.2f cores", float64(metrics.CPURequest-recCPU)/1000.0)
	savingsMem := fmt.Sprintf("%d MB", metrics.MemoryRequest-recMem)

	return models.Recommendation{
		Deployment:              metrics.Deployment,
		CurrentCPURequest:       metrics.CPURequest,
		RecommendedCPU:          recCPU,
		CurrentMemoryRequest:    metrics.MemoryRequest,
		RecommendedMemory:       recMem,
		CPUReductionPercent:     cpuReductionPct,
		MemoryReductionPercent:  memReductionPct,
		ConfidenceScore:         math.Round(confidence*100) / 100, // round to 2 decimals
		Severity:                severity,
		EstimatedMonthlySavings: models.EstimatedSavings{CPU: savingsCPU, Memory: savingsMem},
		Reason:                  reason,
		Warnings:                warnings,
	}
}
