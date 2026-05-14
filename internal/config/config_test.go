package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViperUnmarshal_BindsNestedSnakeCaseDefaults(t *testing.T) {
	v := viper.New()
	v.SetDefault("server.port", 8080)
	v.SetDefault("recommendation.cpu_safety_buffer", DefaultCPUSafetyBuffer)
	v.SetDefault("recommendation.memory_safety_buffer", DefaultMemorySafetyBuffer)
	v.SetDefault("recommendation.min_cpu_millicores", DefaultMinCPUMillicores)
	v.SetDefault("recommendation.min_memory_mb", DefaultMinMemoryMB)

	var cfg Config
	require.NoError(t, v.Unmarshal(&cfg))

	assert.Equal(t, 8080, cfg.Server.Port)
	assert.InDelta(t, DefaultCPUSafetyBuffer, cfg.Recommendation.CPUSafetyBuffer, 1e-9)
	assert.InDelta(t, DefaultMemorySafetyBuffer, cfg.Recommendation.MemorySafetyBuffer, 1e-9)
	assert.Equal(t, DefaultMinCPUMillicores, cfg.Recommendation.MinCPUMillicores)
	assert.Equal(t, DefaultMinMemoryMB, cfg.Recommendation.MinMemoryMB)
}

func TestLoad_RespectsEnvOverrides(t *testing.T) {
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("RECOMMENDATION_CPU_SAFETY_BUFFER", "2.0")
	t.Setenv("RECOMMENDATION_MEMORY_SAFETY_BUFFER", "1.8")
	t.Setenv("RECOMMENDATION_MIN_CPU_MILLICORES", "250")
	t.Setenv("RECOMMENDATION_MIN_MEMORY_MB", "256")

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, 9090, cfg.Server.Port)
	assert.InDelta(t, 2.0, cfg.Recommendation.CPUSafetyBuffer, 1e-9)
	assert.InDelta(t, 1.8, cfg.Recommendation.MemorySafetyBuffer, 1e-9)
	assert.Equal(t, 250, cfg.Recommendation.MinCPUMillicores)
	assert.Equal(t, 256, cfg.Recommendation.MinMemoryMB)
}

func TestNormalize_RepairsInvalidCPUBuffer(t *testing.T) {
	c := RecommendationConfig{
		CPUSafetyBuffer:    0,
		MemorySafetyBuffer: DefaultMemorySafetyBuffer,
		MinCPUMillicores:   DefaultMinCPUMillicores,
		MinMemoryMB:        DefaultMinMemoryMB,
	}
	assert.True(t, c.Normalize())
	assert.InDelta(t, DefaultCPUSafetyBuffer, c.CPUSafetyBuffer, 1e-9)
}

func TestNormalize_ClampHighCPUBuffer(t *testing.T) {
	c := RecommendationConfig{
		CPUSafetyBuffer:    10,
		MemorySafetyBuffer: DefaultMemorySafetyBuffer,
		MinCPUMillicores:   DefaultMinCPUMillicores,
		MinMemoryMB:        DefaultMinMemoryMB,
	}
	assert.True(t, c.Normalize())
	assert.InDelta(t, MaxValidCPUSafetyBuffer, c.CPUSafetyBuffer, 1e-9)
}

func TestValidate_RejectsUnnormalizedInvalid(t *testing.T) {
	c := RecommendationConfig{
		CPUSafetyBuffer:    0.5,
		MemorySafetyBuffer: DefaultMemorySafetyBuffer,
		MinCPUMillicores:   DefaultMinCPUMillicores,
		MinMemoryMB:        DefaultMinMemoryMB,
	}
	assert.Error(t, c.Validate())
}
