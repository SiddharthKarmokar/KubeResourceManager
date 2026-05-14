package services

import (
	"context"

	"github.com/SiddharthKarmokar/KubeResourceManager/internal/domain/models"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/engine/recommendation"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/engine/validation"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/logger"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/observability"
)

// OptimizerService is the primary service for processing metrics and generating recommendations.
type OptimizerService struct {
	recommender *recommendation.Recommender
}

// NewOptimizerService creates a new OptimizerService.
func NewOptimizerService(recommender *recommendation.Recommender) *OptimizerService {
	return &OptimizerService{
		recommender: recommender,
	}
}

// AnalyzeBatch analyzes a batch of metrics and returns recommendations.
func (s *OptimizerService) AnalyzeBatch(ctx context.Context, metricsBatch []models.WorkloadMetrics) ([]models.Recommendation, error) {
	log := logger.FromContext(ctx)

	if err := validation.ValidateMetricsBatch(ctx, metricsBatch); err != nil {
		log.Error("validation failed", "error", err)
		return nil, err
	}

	var recommendations []models.Recommendation
	for _, metrics := range metricsBatch {
		rec := s.recommender.GenerateRecommendation(ctx, metrics)
		recommendations = append(recommendations, rec)

		observability.RecommendationCount.WithLabelValues(rec.Severity).Inc()
		log.Info("generated recommendation",
			"deployment", rec.Deployment,
			"severity", rec.Severity,
			"confidence", rec.ConfidenceScore,
		)
	}

	return recommendations, nil
}
