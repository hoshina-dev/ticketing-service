# ticketing-service

Laboratory experiment ticketing service. Tracks experiment tickets, links them to
experiment templates from Experiment Manager, and coordinates downstream workflows.

## Architecture

- **Fiber v2** HTTP server.
- **Layered**: `handler` → `service` → `repository` (GORM/Postgres).
- **OpenAPI**: swag annotations on handlers; spec served at `/docs/*`.

```
cmd/main.go                 # composition root
internal/
  handler/                  # HTTP transport
  service/                  # business logic
  repository/               # persistence
  clients/copium/         # generated Copium API client (oapi-codegen)
docs/                       # generated swagger for this service
sql/                        # migrations
tools/copium-openapi/       # fetch Copium swagger → OpenAPI 3 for codegen
```

## Dependency services

| service | repo | role |
| --- | --- | --- |
| **Copium** | [github.com/hoshina-dev/copium](https://github.com/hoshina-dev/copium) | Communication API — versioned email templates and async dispatch via an outbox worker. Used to notify clients and lab staff when ticket lifecycle events occur. |
| **Experiment Manager** | [github.com/hoshina-dev/experiment-manager](https://github.com/hoshina-dev/experiment-manager) | Experiment templates and experiment context (forms, calculations). Ticket experiment templates reference template IDs from this service. |

### Copium client (codegen)

The HTTP client for Copium lives in `internal/clients/copium/`. It is generated
from Copium's OpenAPI document with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen).

Copium publishes Swagger 2 at `http://localhost:8081/swagger/doc.json` when running
locally (`make run` in the copium repo). The fetch step converts that spec to
OpenAPI 3 (required for query-parameter schemas) before codegen.

Regenerate after Copium API changes:

```sh
# Copium must be reachable (default URL below)
make codegen-copium

# or override the doc URL
make codegen-copium COPUM_OPENAPI_URL=http://copium.mapfox.hoshina.san/swagger/doc.json
```

This runs `tools/copium-openapi` (fetch + convert), then `go generate` to produce
`client.gen.go` from `openapi.json`. Commit both the spec snapshot and generated
client when the Copium contract changes.

Example usage:

```go
client, err := copium.NewClientWithResponses("http://localhost:8081")
if err != nil { ... }

resp, err := client.SendEmailWithResponse(ctx, copium.SendEmailJSONRequestBody{
    TemplateId: templateID,
    UserId:     &userID,
})
```

## Make targets

| target | what it does |
| --- | --- |
| `make run` | `go run cmd/main.go` |
| `make test` | `go test ./...` |
| `make swagger` | regenerate this service's OpenAPI docs |
| `make codegen-copium` | fetch Copium spec and regenerate `internal/clients/copium` |
| `make generate` | `codegen-copium` + `go generate ./...` |
| `make lint` | `golangci-lint run ./...` |
| `make build` | `go build -o bin/ticketing-service ./cmd` |

## Configuration

Copy `.env.example` to `.env`:

| variable | default | description |
| --- | --- | --- |
| `PORT` | `8080` | HTTP listen port |
| `DATA_SOURCE_NAME` | — | Postgres URL (required) |

## API documentation

| URL | UI |
| --- | --- |
| `http://localhost:8080/docs/index.html` | Swagger UI |
| `http://localhost:8080/docs/doc.json` | Raw OpenAPI 2 spec |

Regenerate after handler annotation changes:

```sh
make swagger
```
