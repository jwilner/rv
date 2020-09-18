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
	@goimports -w -local github.com/jwilner/rv \
		./internal/ \
		./cmd/
	@npx prettier --write src/

proto:
	@protoc -I=. \
		--go_out=:pkg --go_opt=paths=source_relative \
		--go-grpc_out=:pkg --go-grpc_opt=paths=source_relative \
		--js_out=import_style=commonjs:src/ \
		--grpc-web_out=import_style=commonjs,mode=grpcwebtext:src/\
		pb/rvapi/*.proto
# third party deps not shipped with runtime
	@protoc --proto_path third_party/googleapis \
		--js_out=import_style=commonjs:src/pb \
		--grpc-web_out=import_style=commonjs,mode=grpcwebtext:src/pb \
		third_party/googleapis/google/rpc/*.proto
