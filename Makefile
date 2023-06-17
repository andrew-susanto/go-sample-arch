.PHONY: build

mod:
	@go mod tidy -v
	@go mod vendor -v
	@go mod tidy -v

build:
	@echo "BUILD sampleapp..."
	@go build -v -o sampleapp app/*.go

run:
	@echo "RUN sampleapp..."
	make build
	@./sampleapp

test:
	@go test -v -race ./...
