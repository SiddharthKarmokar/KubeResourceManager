# API reference

The HTTP API is a thin JSON interface over the same batch optimizer used by the CLI. There is **no authentication** in the default router; run it behind a private network, service mesh, or API gateway if you expose it beyond localhost.

**Base URL (local):** `http://localhost:8080`

## Conventions

- **Content-Type:** `application/json` for `POST /analyze`.
- **Success responses** return **200** with a JSON **array** body.
- **Errors** return **4xx/5xx** with a JSON object shaped as `{"error": {"code": "...", "message": "..."}}`.

## System endpoints

### `GET /healthz`

Liveness probe. Returns **200** with body `OK`.

### `GET /readyz`

Readiness probe. Returns **200** with body `OK` (extend with dependency checks as needed).

### `GET /metrics`

Prometheus exposition format (text).

### `GET /swagger/*`

Interactive **Swagger UI** plus **`doc.json`** (embedded OpenAPI 2.0 from [swag](https://github.com/swaggo/swag)).

- UI entrypoint: **`/swagger/index.html`** (or `/swagger/` which redirects there).
- Spec JSON: **`/swagger/doc.json`** (served from the registered `swagger` package).

Regenerate specs after changing handler annotations:

```bash
make swagger
```

Source of truth for the committed bundle lives under **`swagger/`** (`swagger.json`, `swagger.yaml`, `docs.go`).

---

## `POST /analyze`

Analyzes one or more workloads and returns an array of recommendations.

### Request body

An array of `WorkloadMetrics`:

| Field | Type | Rules |
|-------|------|--------|
| `deployment` | string | Required, non-empty |
| `cpu_request` | int | Required, **> 0** (millicores) |
| `cpu_usage_avg` | int | Required, **≥ 0** (millicores average) |
| `memory_request` | int | Required, **> 0** (MiB) |
| `memory_usage_avg` | int | Required, **≥ 0** (MiB average) |

**Example**

```json
[
  {
    "deployment": "api-service",
    "cpu_request": 1000,
    "cpu_usage_avg": 180,
    "memory_request": 2048,
    "memory_usage_avg": 700
  }
]
```

### Validation rules (business)

- `cpu_usage_avg` must be **≤** `cpu_request`.
- `memory_usage_avg` must be **≤** `memory_request`.

Violations produce **400** with `INVALID_INPUT`.

### Success response

**200 OK** — JSON array of `Recommendation` objects:

| Field | Type | Description |
|-------|------|----------------|
| `deployment` | string | Workload name |
| `current_cpu_request` | int | Echo of input CPU request |
| `recommended_cpu` | int | Suggested CPU request (millicores) |
| `current_memory_request` | int | Echo of input memory request |
| `recommended_memory` | int | Suggested memory request (MiB) |
| `cpu_reduction_percent` | int | Rounded percent reduction vs current CPU |
| `memory_reduction_percent` | int | Rounded percent reduction vs current memory |
| `confidence_score` | float | Heuristic confidence in `[0.10, 0.99]` (or fixed low when degraded) |
| `severity` | string | `low`, `moderate`, `high`, `critical`, or `degraded` |
| `estimated_monthly_savings` | object | Opaque strings for CPU/memory “savings” presentation |
| `reason` | string | Human-readable summary |
| `warnings` | string[] | Includes active buffer factors and data caveats |

**Example (truncated)**

```json
[
  {
    "deployment": "api-service",
    "current_cpu_request": 1000,
    "recommended_cpu": 270,
    "current_memory_request": 2048,
    "recommended_memory": 910,
    "cpu_reduction_percent": 73,
    "memory_reduction_percent": 56,
    "confidence_score": 0.85,
    "severity": "high",
    "estimated_monthly_savings": {
      "cpu": "0.73 cores",
      "memory": "1138 MB"
    },
    "reason": "Average usage significantly below requested resources with stable utilization patterns.",
    "warnings": [
      "Recommendation includes 1.5x CPU and 1.3x Memory safety buffer.",
      "Current implementation assumes point-in-time average utilization data only."
    ]
  }
]
```

### Error responses

| Status | When |
|--------|------|
| **400** | Malformed JSON, validation failure, or empty batch |
| **500** | Unexpected internal failure |

**Shape**

```json
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "CPU usage cannot exceed CPU request for this analysis"
  }
}
```

**Codes** (non-exhaustive): `INVALID_INPUT`, `INVALID_REQUEST`, `INTERNAL_ERROR`, `NOT_FOUND`.

### `curl` example

```bash
curl -sS -X POST http://localhost:8080/analyze \
  -H "Content-Type: application/json" \
  -d '[{
    "deployment": "api-service",
    "cpu_request": 1000,
    "cpu_usage_avg": 180,
    "memory_request": 2048,
    "memory_usage_avg": 700
  }]'
```

## Design note: averages only

The API intentionally operates on **averages** (not P99/P999). Confidence and warnings call this out. Treat outputs as **governance hints**, not guaranteed SLO-safe limits, unless you enrich the model with tail metrics.
