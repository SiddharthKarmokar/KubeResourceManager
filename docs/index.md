# Kubernetes Resource Optimization Agent

Welcome to the documentation for the **Kubernetes Resource Optimization Agent** — a Go service and CLI that turns **point-in-time average utilization** and **current requests** into **conservative** CPU and memory request recommendations.

## Who is this for?

- Platform engineers operating Kubernetes clusters who want **actionable right-sizing hints** without replacing VPA or building a full metrics warehouse.
- Teams that already export **average CPU/memory usage** per workload and want a **bounded, explainable** recommendation payload.

## What this project does

- Accepts a **batch of workload metrics** (JSON) via **HTTP** or the **CLI**.
- Applies **configurable safety buffers**, **minimum floors**, and **post-processing guardrails** so recommendations stay credible under bad config or sparse metrics.
- Exposes **Prometheus metrics**, **health/readiness** endpoints, and optional **Swagger** docs from the API binary.

## Quick links

- **Source code:** [github.com/SiddharthKarmokar/KubeResourceManager](https://github.com/SiddharthKarmokar/KubeResourceManager)
- **Author:** Siddharh Karmokar — [Portfolio](https://siddkarmokar-portfolio.vercel.app/)
- **Live docs (GitHub Pages):** after enabling Pages from the `gh-pages` branch, the site is published at the URL configured in `mkdocs.yml` (`site_url`).

## Documentation map

| Section | What you will learn |
|--------|----------------------|
| [Architecture](architecture.md) | Layers, packages, and design boundaries |
| [API](api.md) | HTTP routes, payloads, and error shape |
| [Recommendation engine](recommendation-engine.md) | Heuristics, scoring, degraded mode |
| [Configuration](configuration.md) | Viper keys, env vars, validation |
| [Local development](local-development.md) | Build, test, lint, MkDocs |
| [Deployment](deployment.md) | Docker Compose and Kubernetes manifests |
| [Observability & scaling](scaling.md) | Metrics, logging, horizontal scaling notes |
| [Future improvements](future-improvements.md) | Roadmap and non-goals |

## Non-goals (today)

These are intentionally **out of scope** for the current codebase:

- Replacing **VPA** or implementing in-cluster controllers that mutate workloads automatically.
- Long-term metrics storage, multi-cluster federation, or streaming analysis.
- JWT auth on the public analyze API (the shipped router is suitable for private networks or API gateways in front).

Use the sidebar to continue.
