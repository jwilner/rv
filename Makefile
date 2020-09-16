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

proto:
	@protoc -I=. \
		pb/rvapi/*.proto \
		--go_out=:pkg --go_opt=paths=source_relative \
		--go-grpc_out=:pkg --go-grpc_opt=paths=source_relative \
		--js_out=import_style=commonjs:src/ \
		--grpc-web_out=import_style=commonjs,mode=grpcwebtext:src/\
