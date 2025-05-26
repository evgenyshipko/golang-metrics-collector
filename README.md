# Metrics Collector (Golang)

Project for learning golang. 

Agent collects metrics and send it to server. Server store metric values in memory.

### Get started


0. Start PostgreSQL database

```bash
docker run -d --name metrics-collector-pg -p 5433:5432 -e POSTGRES_PASSWORD=metrics -e POSTGRES_USER=metrics -e POSTGRES_DB=metrics postgres
```

1. Install [goose](https://github.com/pressly/goose?tab=readme-ov-file#up)

(For MacOs):
```bash
brew install goose
```

2. Run migrations
```bash
goose -dir internal/server/db/migrations postgres "postgres://metrics:metrics@localhost:5433/metrics?sslmode=disable" up 
```

3. Run metrics collector (server)

```bash
go run ./cmd/server -d="postgres://metrics:metrics@localhost:5433/metrics?sslmode=disable" -m=false
```

4. Run agent (source of metric values)

```bash
go run ./cmd/agent
```

### Unit tests

To run unit test execute command:
```bash
go test ./... -v
```

### Run multichecker

1. Build checker
```bash
go build -o checker cmd/checker/main.go
```

2. Run check
```bash
./checker ./...
```
