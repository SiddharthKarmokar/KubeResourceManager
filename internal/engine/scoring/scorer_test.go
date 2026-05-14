package scoring

import (
	"context"
	"testing"

	"github.com/SiddharthKarmokar/KubeResourceManager/internal/domain/enums"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestCalculateConfidence_DegradedReturnsLow(t *testing.T) {
	s := NewScorer()
	ctx := context.Background()
	m := models.WorkloadMetrics{CPURequest: 1000, CPUUsageAvg: 180, MemoryRequest: 2048, MemoryUsageAvg: 700}
	conf := s.CalculateConfidence(ctx, m, 270, 910, true)
	assert.InDelta(t, 0.25, conf, 1e-9)
}

func TestCalculateConfidence_NoDivideByZeroWhenUsageZero(t *testing.T) {
	s := NewScorer()
	ctx := context.Background()
	m := models.WorkloadMetrics{CPURequest: 100, CPUUsageAvg: 0, MemoryRequest: 200, MemoryUsageAvg: 100}
	conf := s.CalculateConfidence(ctx, m, 50, 150, false)
	assert.GreaterOrEqual(t, conf, 0.10)
	assert.LessOrEqual(t, conf, 0.99)
}

func TestCalculateSeverity_DegradedOverridesHighReduction(t *testing.T) {
	s := NewScorer()
	ctx := context.Background()
	sev := s.CalculateSeverity(ctx, 100, 100, true)
	assert.Equal(t, enums.SeverityDegraded.String(), sev)
}

func TestCalculateSeverity_WastePercentWhenHealthy(t *testing.T) {
	s := NewScorer()
	ctx := context.Background()
	assert.Equal(t, enums.SeverityCritical.String(), s.CalculateSeverity(ctx, 85, 10, false))
	assert.Equal(t, enums.SeverityLow.String(), s.CalculateSeverity(ctx, 5, 5, false))
}
