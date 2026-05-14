# Local development

## Prerequisites

- **Go** (see `go.mod` for the toolchain line)
- **Make** (optional but recommended)
- **Docker** / **Docker Compose** (for container workflows)
- **Python 3** + **pip** (for MkDocs)
- **golangci-lint** (optional; `make lint` skips if missing)
- **swag** (optional; for regenerating Swagger)

## Common commands

```bash
# Format
gofmt -w .

# Unit tests
go test ./...

# Same as above via Makefile (verbose + cover)
make test

# Static analysis
make lint        # go vet + golangci-lint (if installed)

# Build binaries to bin/
make build

# Run API on :8080 (after setting env if needed)
make run
```

## CLI workflow

```bash
make build
./bin/k8s-opt analyze --file path/to/metrics.json
```

The CLI reads the same JSON array schema as `POST /analyze`.

## Documentation (MkDocs)

Install Python dependencies and build the static site:

```bash
pip install -r docs/requirements.txt
mkdocs serve    # local preview
mkdocs build    # outputs to ./site
```

Or use **`make docs`** (installs nothing automatically — ensure `mkdocs` is on `PATH`).

## Swagger generation

If you change handler annotations:

```bash
make swagger
```

## Cross-platform notes

- `Makefile` targets use Unix-style helpers; on Windows, **Git Bash**, **WSL**, or running `go` commands directly works well.
