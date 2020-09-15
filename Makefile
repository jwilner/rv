.PHONY: test migrate dbup

test:
	@go test ./...

gen: migrate
	@go run ./cmd/genmodels

migrate: dbup
	@go run ./cmd/migrate

dbup:
	@docker-compose up -d pg

format:
	@goimports -w -local github.com/jwilner/rv .
