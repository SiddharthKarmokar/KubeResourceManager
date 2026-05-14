package heuristics

import (
	"context"
	"math"

	"github.com/siddk/kube-resource-manager/internal/config"
)

// CalculateSafeCPU computes the recommended CPU request based on heuristics.
func CalculateSafeCPU(ctx context.Context, avgUsage int, cfg config.RecommendationConfig) int {
	buf := math.Max(cfg.CPUSafetyBuffer, config.MinValidCPUSafetyBuffer)
	recommended := int(math.Round(float64(avgUsage) * buf))

	minV := cfg.MinCPUMillicores
	if minV < 1 {
		minV = 1
	}
	if recommended < minV {
		recommended = minV
	}

	return recommended
}

// CalculateSafeMemory computes the recommended Memory request based on conservative memory safety heuristics.
func CalculateSafeMemory(ctx context.Context, avgUsage int, cfg config.RecommendationConfig) int {
	buf := math.Max(cfg.MemorySafetyBuffer, config.MinValidMemorySafetyBuffer)
	recommended := int(math.Round(float64(avgUsage) * buf))

	minV := cfg.MinMemoryMB
	if minV < 1 {
		minV = 1
	}
	if recommended < minV {
		recommended = minV
	}

	return recommended
}
