# Metrics Collector (Golang)

Project for learning golang. 

Agent collects metrics and send it to server. Server store metric values in memory.

### Get started


0. Start PostgreSQL database

```bash
docker run -d --name metrics-collector-pg -p 5433:5432 -e POSTGRES_PASSWORD=metrics -e POSTGRES_USER=metrics -e POSTGRES_DB=metrics postgres
```

1. Run metrics collector (server)

```bash
go run ./cmd/server
```

2. Run agent (source of metric values)

```bash
go run ./cmd/agent
```

### Unit tests

To run unit test execute command:
```bash
go test ./... -v
```
