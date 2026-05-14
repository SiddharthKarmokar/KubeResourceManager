# Architecture

The agent follows **hexagonal (ports & adapters) style** without over-building abstractions: domain types live in one place, HTTP is a thin adapter, and orchestration sits in a small service layer.

## Package layout

| Path | Responsibility |
|------|----------------|
| `cmd/api` | HTTP server entrypoint: config load, wiring, graceful shutdown |
| `cmd/cli` | Cobra CLI: read JSON file, call the same optimizer service |
| `internal/api` | Router, handlers, middleware, DTOs |
| `internal/config` | Viper load, **mapstructure-tagged** structs, **normalize + validate** |
| `internal/domain` | Models (`WorkloadMetrics`, `Recommendation`) and enums (`Severity*`) |
| `internal/engine/heuristics` | Usage × buffer math with floors |
| `internal/engine/recommendation` | Orchestrates heuristics, caps, **finalize** guardrails |
| `internal/engine/scoring` | Confidence and severity (including **degraded** path) |
| `internal/engine/validation` | Request validation (struct tags + business rules) |
| `internal/services` | `OptimizerService` batch loop + observability hooks |
| `internal/logger` | `slog` setup and context logger |
| `internal/observability` | Prometheus metric vectors |
| `deployments/` | Docker / Kubernetes samples |

There is **no `pkg/`** export layer: everything meaningful stays under `internal/` by design.

## Request flow (HTTP)

1. **Handler** decodes JSON into `dto.AnalyzeRequest` (`[]WorkloadMetrics`).
2. **Validation** rejects empty batches and inconsistent metrics (for example usage above request — see `internal/engine/validation`).
3. **Recommender** runs per workload: **heuristics** → **cap to current request** → **finalize invariants** → **scoring**.
4. **Response** is a JSON array of `Recommendation` objects (not wrapped in an envelope).

```text
Client → POST /analyze → Validate batch → Recommender (per item) → JSON []
```

## Configuration safety

Configuration is loaded once at startup. Invalid floats (NaN/Inf), buffers outside `[1.0, 3.0]`, or minimums below platform floors trigger **normalization** (with a structured warning) and then **validation**. This prevents silent “multiply by zero” failures that would otherwise collapse recommendations.

## Observability

- Counters and histograms for HTTP traffic (`internal/observability`).
- Recommendation counter labeled by **severity** (includes `degraded` when finalize repair occurs).

## Extension points

- Swap the heuristic functions for percentile-based models while keeping **finalize** as the last line of defense.
- Add persistence or auth **outside** `internal/engine` so the core remains testable without a database.
