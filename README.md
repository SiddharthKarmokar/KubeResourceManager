# Kubernetes Resource Optimization Agent

A production-grade Kubernetes Resource Optimization Agent built in Golang. This system analyzes Kubernetes workload resource usage and generates safe CPU and memory optimization recommendations based on defined heuristics.

**Author:** [Siddharh Karmokar](https://siddkarmokar-portfolio.vercel.app/) · [Portfolio](https://siddkarmokar-portfolio.vercel.app/)

## Features
- **REST API**: Analyzes workload metrics and returns resource recommendations.
- **CLI**: A local command-line interface for batch processing JSON files.
- **Conservative Heuristics**: Avoids aggressive downsizing and applies mathematically sound safety buffers.
- **Production Ready**: Includes graceful shutdown, structured logging, Prometheus metrics, and containerization.
- **Clean Architecture**: Designed without over-engineering, using Hexagonal Architecture principles.

## Non-Functional Goals
- Real-time streaming metrics.
- Full VPA (Vertical Pod Autoscaler) replacement.
- Mutating controllers or autonomous workload patching.
- Multi-cluster synchronization or long-term metrics storage.

## Quickstart

### Prerequisites
- Go 1.25+
- Docker & Docker Compose
- Make

### Running Locally
```bash
# Build the project
make build

# Start the API via Docker Compose
docker-compose up -d

# Run the CLI
./bin/k8s-opt analyze --file sample.json
```

### API Usage
```bash
curl -X POST http://localhost:8080/analyze \
  -H "Content-Type: application/json" \
  -d '[{
    "deployment": "api-service",
    "cpu_request": 1000,
    "cpu_usage_avg": 180,
    "memory_request": 2048,
    "memory_usage_avg": 700
  }]'
```
### Example Output
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
    },
    {
        "deployment": "worker-service",
        "current_cpu_request": 500,
        "recommended_cpu": 500,
        "current_memory_request": 1024,
        "recommended_memory": 1024,
        "cpu_reduction_percent": 0,
        "memory_reduction_percent": 0,
        "confidence_score": 0.65,
        "severity": "low",
        "estimated_monthly_savings": {
            "cpu": "0.00 cores",
            "memory": "0 MB"
        },
        "reason": "Usage is close to requested resources; no change recommended.",
        "warnings": [
            "Recommendation includes 1.5x CPU and 1.3x Memory safety buffer.",
            "Current implementation assumes point-in-time average utilization data only."
        ]
    }
]
```


## Documentation
Full documentation is available via MkDocs:
```bash
make docs
```

**Published docs (GitHub Pages):** [https://siddk.github.io/kube-resource-manager/](https://siddk.github.io/kube-resource-manager/)

After forking, enable Pages from the `gh-pages` branch and update `site_url` in `mkdocs.yml` if your GitHub Pages URL differs.

### OpenAPI / Swagger UI
Generate the OpenAPI bundle (writes to `swagger/`), then run the API:

```bash
make swagger
go run ./cmd/api
```

Open **`http://localhost:8080/swagger/index.html`**. The UI loads **`/swagger/doc.json`**, which is registered from the generated `swagger` package (`cmd/api` imports it with a blank import).
