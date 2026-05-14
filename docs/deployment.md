# Deployment

## Docker Compose (local)

The repository includes `docker-compose.yml` for a single-service deployment. Typical flow:

```bash
docker compose up -d --build
```

The API listens on the host port mapped in the compose file (defaults align with `SERVER_PORT` / `8080`).

## Container image

Build manually:

```bash
make docker-build
# or
docker build -t k8s-opt:latest .
```

The `Dockerfile` multi-stage build produces a minimal runtime image with the API binary.

## Kubernetes

Manifests live under `deployments/kubernetes/`:

- `deployment.yaml` — Deployment + Service skeleton
- `configmap.yaml` — Example recommendation and server settings

Wire the ConfigMap into the Deployment via envFrom or individual `env` entries, and set resource requests/limits appropriate to your cluster policy.

## GitHub Pages (documentation site)

CI builds MkDocs on pushes to `main` (when docs change) and publishes the static `site/` output to the **`gh-pages`** branch.

**Enable Pages:** Repository **Settings → Pages → Build and deployment → Source: Deploy from a branch → `gh-pages` / `/ (root)`**.

Canonical URLs are set in `mkdocs.yml` for **`SiddharthKarmokar/KubeResourceManager`**; adjust `site_url` / `repo_url` if you rename the repository.

## Production hardening checklist

- Place the API **behind authentication** (gateway, mesh, or private networking).
- Scrape **`/metrics`** with Prometheus and alert on error rates and latency.
- Run **`/readyz`** checks that validate any future external dependencies.
- Pin container images by digest in Kubernetes; avoid `:latest` in prod.
