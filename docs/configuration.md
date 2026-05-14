# Configuration

Configuration is loaded at process start via **Viper** (`internal/config`). Values can come from **defaults**, **environment variables**, or (if you extend `Load`) config files — the shipped `Load()` focuses on env + defaults.

## Environment variables

Viper maps nested keys with dots to env names using **`_`** (see `SetEnvKeyReplacer`).

| Variable | Maps to | Default | Notes |
|----------|---------|---------|--------|
| `SERVER_PORT` | `server.port` | `8080` | HTTP listen port |
| `RECOMMENDATION_CPU_SAFETY_BUFFER` | `recommendation.cpu_safety_buffer` | `1.5` | Clamped to **[1.0, 3.0]** |
| `RECOMMENDATION_MEMORY_SAFETY_BUFFER` | `recommendation.memory_safety_buffer` | `1.3` | Clamped to **[1.0, 3.0]** |
| `RECOMMENDATION_MIN_CPU_MILLICORES` | `recommendation.min_cpu_millicores` | `100` | Floors below **10** repair to default |
| `RECOMMENDATION_MIN_MEMORY_MB` | `recommendation.min_memory_mb` | `128` | Floors below **32** repair to default |

### Kubernetes example

See `deployments/kubernetes/configmap.yaml` for a ConfigMap that mirrors the variables above.

## Struct binding (`mapstructure`)

Nested configuration uses **explicit `mapstructure` tags** so snake_case keys bind reliably into Go structs. This avoids the classic pitfall where `Unmarshal` succeeds but leaves **zero values** for critical floats and ints.

## Normalize + validate

On startup:

1. **Normalize** repairs NaN/Inf, sub-minimum buffers, above-max buffers (high values clamp to the max instead of silently accepting runaway multipliers), and minimums that are too small for practical clusters.
2. A **warning** is logged if anything was repaired.
3. **Validate** enforces invariants; failure prevents the process from serving traffic.

## Operational guidance

- Treat buffers as **organizational policy**: higher buffers are more conservative; values outside `[1, 3]` are rejected or clamped.
- Minimums should reflect your **platform baseline** (for example smallest workload you are willing to schedule).
- For GitOps, prefer pinning these values in a **ConfigMap** or external secrets manager rather than ad-hoc shell exports.
