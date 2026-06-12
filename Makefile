.PHONY: test run infra-up infra-down

GOCACHE ?= /private/tmp/wingops-go-build

test:
	GOCACHE=$(GOCACHE) go test ./...

run:
	GOCACHE=$(GOCACHE) go run ./cmd/server

infra-up:
	docker compose up -d postgres redis nats victoriametrics

infra-down:
	docker compose down
