package validation

import (
	"context"
	"testing"

	"github.com/siddk/kube-resource-manager/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestValidateWorkloadMetrics(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		metrics models.WorkloadMetrics
		wantErr bool
	}{
		{
			name: "valid metrics",
			metrics: models.WorkloadMetrics{
				Deployment:     "test-app",
				CPURequest:     1000,
				CPUUsageAvg:    500,
				MemoryRequest:  2048,
				MemoryUsageAvg: 1024,
			},
			wantErr: false,
		},
		{
			name: "missing deployment",
			metrics: models.WorkloadMetrics{
				CPURequest:     1000,
				CPUUsageAvg:    500,
				MemoryRequest:  2048,
				MemoryUsageAvg: 1024,
			},
			wantErr: true,
		},
		{
			name: "usage exceeds request",
			metrics: models.WorkloadMetrics{
				Deployment:     "test-app",
				CPURequest:     1000,
				CPUUsageAvg:    1500,
				MemoryRequest:  2048,
				MemoryUsageAvg: 1024,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWorkloadMetrics(ctx, tt.metrics)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
