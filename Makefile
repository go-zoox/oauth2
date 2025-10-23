# Makefile for acme-go

.PHONY: build test clean install run-example

# Build the application
build:
	go build -o bin/acme-go main.go

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/acme-go-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/acme-go-darwin-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/acme-go-windows-amd64.exe main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install to system
install: build
	sudo cp bin/acme-go /usr/local/bin/
	sudo chmod +x /usr/local/bin/acme-go

# Run example commands
run-example:
	@echo "Building acme-go..."
	@go build -o bin/acme-go main.go
	@echo "\nShowing help:"
	@./bin/acme-go --help
	@echo "\nShowing issue command help:"
	@./bin/acme-go issue --help

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Tidy dependencies
tidy:
	go mod tidy

# Development setup
dev-setup:
	go mod download
	@echo "Development environment ready!"