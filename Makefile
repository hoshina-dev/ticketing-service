.PHONY: install run test deps lint generate codegen-copium swagger format build

COPUM_OPENAPI_URL ?= http://localhost:8081/swagger/doc.json

install:
	go mod download

build:
	go build -o bin/ticketing-service ./cmd

run:
	go run cmd/main.go

test:
	go test ./...

deps:
	go mod tidy

lint:
	golangci-lint run ./...

codegen-copium:
	go run ./tools/copium-openapi -url $(COPUM_OPENAPI_URL) -o internal/clients/copium/openapi.json
	go generate ./internal/clients/copium/...

generate: codegen-copium
	go generate ./...

swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/main.go -o docs --parseDependency --parseInternal
	@if grep -q "github_com" docs/swagger.json; then \
		echo "\n⚠️ WARNING: Found 'github_com' in swagger.json!"; \
		echo "Some models are not properly named with @name tags."; \
		echo "Please add @name tags to all model structs in your code."; \
		exit 1; \
	else \
		echo "\n✅ Swagger validation passed: No 'github_com' references found"; \
	fi

format:
	go run github.com/swaggo/swag/cmd/swag@latest fmt
	go fmt ./...
	gofmt -s -w .

.DEFAULT_GOAL = run