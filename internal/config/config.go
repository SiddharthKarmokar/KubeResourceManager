package config

import (
	"fmt"
	"log/slog"
	"math"
	"strings"

	"github.com/spf13/viper"
)

// Safe defaults and bounds for recommendation heuristics.
const (
	DefaultCPUSafetyBuffer    = 1.5
	DefaultMemorySafetyBuffer = 1.3
	DefaultMinCPUMillicores   = 100
	DefaultMinMemoryMB        = 128

	MinValidCPUSafetyBuffer    = 1.0
	MaxValidCPUSafetyBuffer    = 3.0
	MinValidMemorySafetyBuffer = 1.0
	MaxValidMemorySafetyBuffer = 3.0

	MinAllowedCPUMillicores = 10
	MinAllowedMemoryMB      = 32
)

// Config represents the application configuration.
type Config struct {
	Server         ServerConfig         `mapstructure:"server"`
	Recommendation RecommendationConfig `mapstructure:"recommendation"`
}

// ServerConfig holds the HTTP server settings.
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// RecommendationConfig holds the heuristics and safety buffer settings.
type RecommendationConfig struct {
	CPUSafetyBuffer    float64 `mapstructure:"cpu_safety_buffer"`
	MemorySafetyBuffer float64 `mapstructure:"memory_safety_buffer"`
	MinCPUMillicores   int     `mapstructure:"min_cpu_millicores"`
	MinMemoryMB        int     `mapstructure:"min_memory_mb"`
}

// Normalize repairs invalid recommendation fields using safe defaults and clamps
// out-of-range values. It returns true if any field was changed.
func (c *RecommendationConfig) Normalize() bool {
	repaired := false

	switch {
	case math.IsNaN(c.CPUSafetyBuffer) || math.IsInf(c.CPUSafetyBuffer, 0) || c.CPUSafetyBuffer < MinValidCPUSafetyBuffer:
		c.CPUSafetyBuffer = DefaultCPUSafetyBuffer
		repaired = true
	case c.CPUSafetyBuffer > MaxValidCPUSafetyBuffer:
		c.CPUSafetyBuffer = MaxValidCPUSafetyBuffer
		repaired = true
	}

	switch {
	case math.IsNaN(c.MemorySafetyBuffer) || math.IsInf(c.MemorySafetyBuffer, 0) || c.MemorySafetyBuffer < MinValidMemorySafetyBuffer:
		c.MemorySafetyBuffer = DefaultMemorySafetyBuffer
		repaired = true
	case c.MemorySafetyBuffer > MaxValidMemorySafetyBuffer:
		c.MemorySafetyBuffer = MaxValidMemorySafetyBuffer
		repaired = true
	}

	if c.MinCPUMillicores < MinAllowedCPUMillicores {
		c.MinCPUMillicores = DefaultMinCPUMillicores
		repaired = true
	}

	if c.MinMemoryMB < MinAllowedMemoryMB {
		c.MinMemoryMB = DefaultMinMemoryMB
		repaired = true
	}

	return repaired
}

// Validate returns an error if recommendation configuration violates invariants.
// Call after Normalize().
func (c RecommendationConfig) Validate() error {
	if c.CPUSafetyBuffer < MinValidCPUSafetyBuffer || c.CPUSafetyBuffer > MaxValidCPUSafetyBuffer {
		return fmt.Errorf("recommendation.cpu_safety_buffer out of range [%v, %v]: %v",
			MinValidCPUSafetyBuffer, MaxValidCPUSafetyBuffer, c.CPUSafetyBuffer)
	}
	if c.MemorySafetyBuffer < MinValidMemorySafetyBuffer || c.MemorySafetyBuffer > MaxValidMemorySafetyBuffer {
		return fmt.Errorf("recommendation.memory_safety_buffer out of range [%v, %v]: %v",
			MinValidMemorySafetyBuffer, MaxValidMemorySafetyBuffer, c.MemorySafetyBuffer)
	}
	if c.MinCPUMillicores < MinAllowedCPUMillicores {
		return fmt.Errorf("recommendation.min_cpu_millicores below minimum %d: %d", MinAllowedCPUMillicores, c.MinCPUMillicores)
	}
	if c.MinMemoryMB < MinAllowedMemoryMB {
		return fmt.Errorf("recommendation.min_memory_mb below minimum %d: %d", MinAllowedMemoryMB, c.MinMemoryMB)
	}
	return nil
}

// Load loads the configuration from environment variables or config file.
func Load() (*Config, error) {
	v := viper.New()

	v.SetDefault("server.port", 8080)
	v.SetDefault("recommendation.cpu_safety_buffer", DefaultCPUSafetyBuffer)
	v.SetDefault("recommendation.memory_safety_buffer", DefaultMemorySafetyBuffer)
	v.SetDefault("recommendation.min_cpu_millicores", DefaultMinCPUMillicores)
	v.SetDefault("recommendation.min_memory_mb", DefaultMinMemoryMB)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if cfg.Server.Port <= 0 {
		cfg.Server.Port = 8080
	}

	if repaired := cfg.Recommendation.Normalize(); repaired {
		slog.Default().Warn("recommendation config contained invalid values; safe defaults and bounds were applied")
	}

	if err := cfg.Recommendation.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
