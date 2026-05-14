package recommendation

import (
	"context"
	"testing"

	"github.com/SiddharthKarmokar/KubeResourceManager/internal/config"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/domain/models"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/engine/scoring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRecommendation_ApiServiceScenario(t *testing.T) {
	cfg := config.RecommendationConfig{
		CPUSafetyBuffer:    1.5,
		MemorySafetyBuffer: 1.3,
		MinCPUMillicores:   100,
		MinMemoryMB:        128,
	}
	r := NewRecommender(cfg, scoring.NewScorer())
	m := models.WorkloadMetrics{
		Deployment:     "api-service",
		CPURequest:     1000,
		CPUUsageAvg:    180,
		MemoryRequest:  2048,
		MemoryUsageAvg: 700,
	}

	rec := r.GenerateRecommendation(context.Background(), m)

	assert.Positive(t, rec.RecommendedCPU)
	assert.Positive(t, rec.RecommendedMemory)
	assert.GreaterOrEqual(t, rec.RecommendedCPU, m.CPUUsageAvg)
	assert.GreaterOrEqual(t, rec.RecommendedMemory, m.MemoryUsageAvg)
	assert.Less(t, rec.CPUReductionPercent, 80)
	assert.NotEqual(t, "critical", rec.Severity)
}

func TestGenerateRecommendation_WorkerNearSaturation(t *testing.T) {
	cfg := config.RecommendationConfig{
		CPUSafetyBuffer:    1.5,
		MemorySafetyBuffer: 1.3,
		MinCPUMillicores:   100,
		MinMemoryMB:        128,
	}
	r := NewRecommender(cfg, scoring.NewScorer())
	m := models.WorkloadMetrics{
		Deployment:     "worker-service",
		CPURequest:     500,
		CPUUsageAvg:    450,
		MemoryRequest:  1024,
		MemoryUsageAvg: 900,
	}

	rec := r.GenerateRecommendation(context.Background(), m)

	assert.Equal(t, 500, rec.RecommendedCPU)
	assert.Equal(t, 1024, rec.RecommendedMemory)
	assert.Equal(t, 0, rec.CPUReductionPercent)
	assert.Equal(t, 0, rec.MemoryReductionPercent)
}

func TestGenerateRecommendation_NeverZeroWithBrokenMins(t *testing.T) {
	cfg := config.RecommendationConfig{
		CPUSafetyBuffer:    1.5,
		MemorySafetyBuffer: 1.3,
		MinCPUMillicores:   0,
		MinMemoryMB:        0,
	}
	r := NewRecommender(cfg, scoring.NewScorer())
	m := models.WorkloadMetrics{
		Deployment:     "api-service",
		CPURequest:     1000,
		CPUUsageAvg:    180,
		MemoryRequest:  2048,
		MemoryUsageAvg: 700,
	}

	rec := r.GenerateRecommendation(context.Background(), m)

	require.Positive(t, rec.RecommendedCPU)
	require.Positive(t, rec.RecommendedMemory)
}

func TestFinalizeRecommendations_MarksDegradedWhenDraftBelowUsage(t *testing.T) {
	cfg := config.RecommendationConfig{
		MinCPUMillicores: 100,
		MinMemoryMB:      128,
	}
	m := models.WorkloadMetrics{
		CPURequest:     1000,
		CPUUsageAvg:    500,
		MemoryRequest:  2048,
		MemoryUsageAvg: 700,
	}
	cpu, mem, degraded := finalizeRecommendations(100, 700, m, cfg)
	assert.True(t, degraded)
	assert.Equal(t, 500, cpu)
	assert.Equal(t, 700, mem)
}
