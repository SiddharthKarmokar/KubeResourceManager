package scoring

import (
	"context"

	"github.com/SiddharthKarmokar/KubeResourceManager/internal/domain/enums"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/domain/models"
)

// Scorer calculates confidence and severity.
type Scorer struct{}

// NewScorer creates a new Scorer.
func NewScorer() *Scorer {
	return &Scorer{}
}

// CalculateConfidence estimates a heuristic confidence based on the delta between the new recommendation and the average usage.
// Since only average utilization metrics are available, confidence is heuristic-based rather than statistically predictive.
// When degraded is true, the engine applied emergency clamps; confidence is lowered to avoid over-trusting the result.
func (s *Scorer) CalculateConfidence(ctx context.Context, metrics models.WorkloadMetrics, recCPU, recMem int, degraded bool) float64 {
	if degraded {
		return 0.25
	}

	// In a real system, this would use historical percentiles.
	// For point-in-time averages, a 1.5x buffer generally gives us ~0.85 confidence.
	confidence := 0.85

	if metrics.CPUUsageAvg > 0 && metrics.MemoryUsageAvg > 0 {
		cpuBufferRatio := float64(recCPU) / float64(metrics.CPUUsageAvg)
		memBufferRatio := float64(recMem) / float64(metrics.MemoryUsageAvg)

		if cpuBufferRatio < 1.2 || memBufferRatio < 1.2 {
			confidence -= 0.20
		} else if cpuBufferRatio > 2.0 && memBufferRatio > 2.0 {
			confidence += 0.10
		}
	} else {
		confidence -= 0.15
	}

	if confidence > 0.99 {
		confidence = 0.99
	}
	if confidence < 0.10 {
		confidence = 0.10
	}

	return confidence
}

// CalculateSeverity determines the severity of the overprovisioning based on waste percent.
// When degraded is true, severity reflects uncertain or repaired inputs rather than waste magnitude.
func (s *Scorer) CalculateSeverity(ctx context.Context, cpuReductionPct, memReductionPct int, degraded bool) string {
	if degraded {
		return enums.SeverityDegraded.String()
	}

	maxReduction := cpuReductionPct
	if memReductionPct > maxReduction {
		maxReduction = memReductionPct
	}

	if maxReduction >= 80 {
		return enums.SeverityCritical.String()
	} else if maxReduction >= 50 {
		return enums.SeverityHigh.String()
	} else if maxReduction >= 20 {
		return enums.SeverityModerate.String()
	}

	return enums.SeverityLow.String()
}
