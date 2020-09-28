.PHONY: build

build:
	@echo "BUILD sampleapp..."
	@go build -v -o sampleapp app/*.go

run:
	@echo "RUN sampleapp..."
	make build
	@./sampleapp -port 9000
