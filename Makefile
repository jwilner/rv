.PHONY: test

test:
	go test ./...

migrate:
	go run ./cmd/migrate
