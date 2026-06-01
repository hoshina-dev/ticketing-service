.PHONY: install run test deps lint generate swagger format build

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

generate:
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
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 fmt
	go fmt ./...
	gofmt -s -w .

build:
	go build -o bin/custapi cmd/main.go

seed:
	go run cmd/seed/main.go

.DEFAULT_GOAL = run