.PHONY: test run

GOCACHE ?= /private/tmp/wingops-go-build

test:
	GOCACHE=$(GOCACHE) go test ./...

run:
	GOCACHE=$(GOCACHE) go run ./cmd/server
