include .env

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

format-gci:
	$(LOCAL_BIN)/gci write .

list-needed-format:
	$(LOCAL_BIN)/gci list .

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/daixiang0/gci@v0.13.5
	GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@v3.4.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.23.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.1.0

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make generate-auth-api

generate-auth-api:
	mkdir -p pkg/authv1
	protoc --proto_path api/proto/auth/v1 --proto_path vendor.protogen \
	--go_out=pkg/authv1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/authv1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--grpc-gateway_out=pkg/authv1 --grpc-gateway_opt=paths=source_relative \
    --plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
    --validate_out lang=go:pkg/authv1 --validate_opt=paths=source_relative \
    --plugin=protoc-gen-validate=bin/protoc-gen-validate \
	api/proto/auth/v1/auth_service.proto \
	api/proto/auth/v1/auth_user_roles.proto

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

create-migration:
	$(LOCAL_BIN)/goose -dir migrations create add_user_table sql

clean-mod-cache:
	GOBIN=$(LOCAL_BIN) go clean -

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/waryataw/auth/internal/service/...,github.com/waryataw/auth/internal/api/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore

vendor-proto:
		@if [ ! -d vendor.protogen/validate ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
			mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
			git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
			mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
			rm -rf vendor.protogen/openapiv2 ;\
		fi