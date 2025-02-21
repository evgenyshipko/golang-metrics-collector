# Metrics Collector (Golang)

Project for learning golang. 

Agent collects metrics and send it to server. Server store metric values in memory.

### Get started

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
