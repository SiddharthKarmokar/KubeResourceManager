package heuristics

import (
	"context"
	"testing"

	"github.com/siddk/kube-resource-manager/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestCalculateSafeCPU(t *testing.T) {
	cfg := config.RecommendationConfig{
		CPUSafetyBuffer:  1.5,
		MinCPUMillicores: 100,
	}
	ctx := context.Background()

	tests := []struct {
		name     string
		avgUsage int
		expected int
	}{
		{"normal usage", 200, 300},
		{"below min threshold", 50, 100},
		{"zero usage", 0, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateSafeCPU(ctx, tt.avgUsage, cfg)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateSafeMemory(t *testing.T) {
	cfg := config.RecommendationConfig{
		MemorySafetyBuffer: 1.3,
		MinMemoryMB:        128,
	}
	ctx := context.Background()

	tests := []struct {
		name     string
		avgUsage int
		expected int
	}{
		{"normal usage", 1000, 1300},
		{"below min threshold", 50, 128},
		{"zero usage", 0, 128},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateSafeMemory(ctx, tt.avgUsage, cfg)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateSafeCPU_ZeroSafetyBufferStillRespectsFloor(t *testing.T) {
	cfg := config.RecommendationConfig{
		CPUSafetyBuffer:  0,
		MinCPUMillicores: 100,
	}
	ctx := context.Background()
	result := CalculateSafeCPU(ctx, 200, cfg)
	assert.Equal(t, 200, result)
}

func TestCalculateSafeMemory_ZeroSafetyBufferStillRespectsFloor(t *testing.T) {
	cfg := config.RecommendationConfig{
		MemorySafetyBuffer: 0,
		MinMemoryMB:        128,
	}
	ctx := context.Background()
	result := CalculateSafeMemory(ctx, 200, cfg)
	assert.Equal(t, 200, result)
}
