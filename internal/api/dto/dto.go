package dto

import "github.com/siddk/kube-resource-manager/internal/domain/models"

// AnalyzeRequest is the incoming batch of metrics.
type AnalyzeRequest []models.WorkloadMetrics

// AnalyzeResponse is the outgoing batch of recommendations.
type AnalyzeResponse []models.Recommendation
