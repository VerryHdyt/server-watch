# server-watch

A lightweight Linux server monitoring service written in Go. It exposes a small HTTP API for health checks and host resource metrics.

## Features

- `GET /health` health check
- `GET /metrics` JSON metrics
- Linux memory metrics from `/proc/meminfo`
- Root filesystem usage metrics
- Docker image and GitHub Actions CI
- Unit and Linux procfs tests

> CPU sampling, configurable HTTP checks, and Prometheus exposition are planned for later milestones.

## Requirements

- Go 1.22+
- Linux host for real system metrics
- Docker (optional)

## Run locally

```bash
go test ./...
go run .
```

The service listens on `http://localhost:8080`. Set a custom address with:

```bash
SERVER_WATCH_ADDR=:9090 go run .
```

## API

### Health

```bash
curl http://localhost:8080/health
```

Returns `{"status":"ok"}`.

### Metrics

```bash
curl http://localhost:8080/metrics
```

Returns hostname, memory usage, and root filesystem usage as JSON. Linux procfs is required for memory metrics. In Docker, the values describe the container's visible environment unless host resources are explicitly exposed.

## Docker

```bash
docker build -t server-watch .
docker run --rm -p 8080:8080 server-watch
```

## Roadmap

- [ ] CPU usage sampling from `/proc/stat`
- [ ] Configurable HTTP endpoint checks
- [ ] Prometheus-compatible metrics
- [ ] Docker Compose example
- [ ] Integration tests and release workflow

## Project structure

```text
.
├── .github/workflows/ci.yml
├── Dockerfile
├── go.mod
├── main.go
├── metrics.go
└── metrics_test.go
```

## License

No license has been added yet.
