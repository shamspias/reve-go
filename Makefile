.PHONY: all build test coverage lint fmt vet clean tidy doc

all: fmt vet lint test

build:
	go build -v ./...

test:
	go test -v -race ./...

coverage:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage: coverage.html"

lint:
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -f coverage.out coverage.html
	go clean -cache -testcache

tidy:
	go mod tidy

doc:
	@echo "http://localhost:6060/pkg/github.com/shamspias/reve-go/"
	@which godoc > /dev/null || go install golang.org/x/tools/cmd/godoc@latest
	godoc -http=:6060

# Examples
example-basic:
	@test -n "$$REVE_API_KEY" || (echo "REVE_API_KEY required" && exit 1)
	go run examples/basic/main.go

example-create:
	@test -n "$$REVE_API_KEY" || (echo "REVE_API_KEY required" && exit 1)
	go run examples/create/main.go

example-edit:
	@test -n "$$REVE_API_KEY" || (echo "REVE_API_KEY required" && exit 1)
	go run examples/edit/main.go

example-remix:
	@test -n "$$REVE_API_KEY" || (echo "REVE_API_KEY required" && exit 1)
	go run examples/remix/main.go

example-batch:
	@test -n "$$REVE_API_KEY" || (echo "REVE_API_KEY required" && exit 1)
	go run examples/batch/main.go

example-proxy:
	@test -n "$$REVE_API_KEY" || (echo "REVE_API_KEY required" && exit 1)
	go run examples/proxy/main.go

example-errors:
	@test -n "$$REVE_API_KEY" || (echo "REVE_API_KEY required" && exit 1)
	go run examples/error-handling/main.go

example-complete:
	@test -n "$$REVE_API_KEY" || (echo "REVE_API_KEY required" && exit 1)
	go run examples/complete/main.go

bench:
	go test -bench=. -benchmem ./...

security:
	@which govulncheck > /dev/null || go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/godoc@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest

release: fmt vet lint test security
	@echo "Ready for release!"

help:
	@echo "Targets:"
	@echo "  all             - fmt, vet, lint, test"
	@echo "  build           - Build package"
	@echo "  test            - Run tests"
	@echo "  coverage        - Test with coverage"
	@echo "  lint            - Run linter"
	@echo "  fmt             - Format code"
	@echo "  vet             - Run go vet"
	@echo "  clean           - Clean artifacts"
	@echo "  tidy            - Tidy dependencies"
	@echo "  doc             - Start godoc server"
	@echo "  example-*       - Run examples"
	@echo "  bench           - Run benchmarks"
	@echo "  security        - Security check"
	@echo "  tools           - Install dev tools"
	@echo "  release         - Prepare release"
