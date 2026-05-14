package validation

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/siddk/kube-resource-manager/internal/domain/models"
	"github.com/siddk/kube-resource-manager/internal/errors"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateWorkloadMetrics validates a single WorkloadMetrics instance.
func ValidateWorkloadMetrics(ctx context.Context, metrics models.WorkloadMetrics) error {
	if err := validate.Struct(metrics); err != nil {
		return errors.NewAPIError(errors.CodeInvalidInput, err.Error())
	}

	// Additional business logic validation
	if metrics.CPUUsageAvg > metrics.CPURequest {
		return errors.NewAPIError(errors.CodeInvalidInput, "CPU usage cannot exceed CPU request for this analysis")
	}

	if metrics.MemoryUsageAvg > metrics.MemoryRequest {
		return errors.NewAPIError(errors.CodeInvalidInput, "Memory usage cannot exceed Memory request for this analysis")
	}

	return nil
}

// ValidateMetricsBatch validates a batch of WorkloadMetrics.
func ValidateMetricsBatch(ctx context.Context, metrics []models.WorkloadMetrics) error {
	if len(metrics) == 0 {
		return errors.NewAPIError(errors.CodeInvalidInput, "empty metrics array")
	}

	for _, m := range metrics {
		if err := ValidateWorkloadMetrics(ctx, m); err != nil {
			return err
		}
	}

	return nil
}
