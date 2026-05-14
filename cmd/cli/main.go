package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/SiddharthKarmokar/KubeResourceManager/internal/config"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/domain/models"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/engine/recommendation"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/engine/scoring"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/logger"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/services"
	"github.com/spf13/cobra"
)

func main() {
	logger.Init()

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	scorer := scoring.NewScorer()
	recommender := recommendation.NewRecommender(cfg.Recommendation, scorer)
	optimizer := services.NewOptimizerService(recommender)

	var inputFile string

	var analyzeCmd = &cobra.Command{
		Use:   "analyze",
		Short: "Analyze workload metrics",
		Long:  `Analyze Kubernetes workload resource usage and generate safe CPU and memory optimization recommendations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if inputFile == "" {
				return fmt.Errorf("--file is required")
			}

			data, err := os.ReadFile(inputFile)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			var metrics []models.WorkloadMetrics
			if err := json.Unmarshal(data, &metrics); err != nil {
				return fmt.Errorf("failed to parse JSON: %w", err)
			}

			ctx := context.Background()
			recommendations, err := optimizer.AnalyzeBatch(ctx, metrics)
			if err != nil {
				return fmt.Errorf("analysis failed: %w", err)
			}

			output, err := json.MarshalIndent(recommendations, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to format output: %w", err)
			}

			fmt.Println(string(output))
			return nil
		},
	}

	analyzeCmd.Flags().StringVarP(&inputFile, "file", "f", "", "Path to the JSON file containing workload metrics")

	var rootCmd = &cobra.Command{Use: "k8s-opt"}
	rootCmd.AddCommand(analyzeCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
