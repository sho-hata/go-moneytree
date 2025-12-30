.PHONY: test lint tools coverage

test:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

lint:
	golangci-lint run

tools:
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing octocov..."
	@go install github.com/k1LoW/octocov@latest

coverage: test
	octocov

