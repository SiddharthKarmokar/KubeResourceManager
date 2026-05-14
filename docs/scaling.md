# Observability and scaling

## Metrics

The router installs a Prometheus handler at **`GET /metrics`**. Histograms and counters include HTTP method, path, and status labels suitable for SRE dashboards.

**Recommendation volume** is tracked per **`severity`** label (values include `degraded` when finalize guardrails fire). Watch for spikes in `degraded` — they often indicate upstream metric quality issues or misconfigured environment overrides.

## Logging

The API uses **`log/slog`** with JSON output. Request middleware attaches contextual fields; the optimizer logs per-deployment severity and confidence at **info** level.

## Horizontal scaling

The service is **stateless**: you can run multiple replicas behind a load balancer as long as:

- Configuration is **consistent** across pods (same buffers/minimums unless you intentionally shard policy).
- Clients can **POST full metric batches** (no server-side session).

## Performance considerations

- Batch analysis is **O(n)** over workloads with minimal allocations in the hot path; tune **request body limits** at the ingress if untrusted clients can post huge arrays.
- Keep **timeouts** (`internal/api/router`) aligned with your ingress and client SLAs.

## Rate limiting

Not included in-tree. Terminate at **API gateway** (for example Kong, Envoy, NGINX) or service mesh policy when exposing beyond trusted networks.
