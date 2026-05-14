package recommendation

import (
	"github.com/siddk/kube-resource-manager/internal/config"
	"github.com/siddk/kube-resource-manager/internal/domain/models"
)

// finalizeRecommendations enforces conservative invariants: recommendations never
// sit below observed average usage, never zero when the workload has a positive
// request, and never below configured minimums. It returns degraded=true when
// emergency clamps were required (invalid draft or zero draft with active request).
func finalizeRecommendations(draftCPU, draftMem int, m models.WorkloadMetrics, cfg config.RecommendationConfig) (cpu, mem int, degraded bool) {
	cpu = draftCPU
	mem = draftMem

	if m.CPURequest > 0 {
		if cpu == 0 || cpu < m.CPUUsageAvg {
			degraded = true
		}
		cpu = max(cpu, m.CPUUsageAvg)
		cpu = max(cpu, cfg.MinCPUMillicores)
		cpu = max(cpu, 1)
		cpu = min(cpu, m.CPURequest)
	}

	if m.MemoryRequest > 0 {
		if mem == 0 || mem < m.MemoryUsageAvg {
			degraded = true
		}
		mem = max(mem, m.MemoryUsageAvg)
		mem = max(mem, cfg.MinMemoryMB)
		mem = max(mem, 1)
		mem = min(mem, m.MemoryRequest)
	}

	return cpu, mem, degraded
}
