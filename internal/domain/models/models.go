package models

// WorkloadMetrics represents the incoming metrics for a single Kubernetes workload.
// @Description Workload identifier plus CPU/memory requests and average utilization (millicores and MiB).
// @name WorkloadMetrics
type WorkloadMetrics struct {
	Deployment     string `json:"deployment" validate:"required" example:"api-service"`
	CPURequest     int    `json:"cpu_request" validate:"required,gt=0" example:"1000"`
	CPUUsageAvg    int    `json:"cpu_usage_avg" validate:"required,gte=0" example:"180"`
	MemoryRequest  int    `json:"memory_request" validate:"required,gt=0" example:"2048"`
	MemoryUsageAvg int    `json:"memory_usage_avg" validate:"required,gte=0" example:"700"`
}

// EstimatedSavings represents the resource savings.
// @name EstimatedSavings
type EstimatedSavings struct {
	CPU    string `json:"cpu" example:"0.73 cores"`
	Memory string `json:"memory" example:"1138 MB"`
}

// Recommendation represents the suggested resource limits and rationale.
// @name Recommendation
type Recommendation struct {
	Deployment              string           `json:"deployment"`
	CurrentCPURequest       int              `json:"current_cpu_request"`
	RecommendedCPU          int              `json:"recommended_cpu"`
	CurrentMemoryRequest    int              `json:"current_memory_request"`
	RecommendedMemory       int              `json:"recommended_memory"`
	CPUReductionPercent     int              `json:"cpu_reduction_percent"`
	MemoryReductionPercent  int              `json:"memory_reduction_percent"`
	ConfidenceScore         float64          `json:"confidence_score"`
	Severity                string           `json:"severity"`
	EstimatedMonthlySavings EstimatedSavings `json:"estimated_monthly_savings"`
	Reason                  string           `json:"reason"`
	Warnings                []string         `json:"warnings,omitempty"`
}
