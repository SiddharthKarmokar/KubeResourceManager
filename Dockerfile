# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o k8s-opt-api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o k8s-opt-cli ./cmd/cli/main.go

# Run stage
FROM alpine:3.19

# Install curl for healthchecks and ca-certificates for external API calls
RUN apk --no-cache add ca-certificates curl

WORKDIR /app

# Non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Copy binaries and assign ownership to the non-root user
COPY --from=builder --chown=appuser:appgroup /app/k8s-opt-api .
COPY --from=builder --chown=appuser:appgroup /app/k8s-opt-cli .

EXPOSE 8080

# (Note: This will be overridden by docker-compose, but acts as a great fallback)
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/healthz || exit 1

ENTRYPOINT ["./k8s-opt-api"]