PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin
$(shell [ -f bin] || mkdir -p $(PROJECT_BIN))
PATH := $(PROJECT_BIN):$(PATH)

GOLANGCI_LINT = $(PROJECT_BIN)/golanci-lint

.PHONY: .install-linter
.install-linter:
	### INSTALL GOLANGCI_LINT ###
	[ -f $(PROJECT_BIN)/golanci-lint ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) v1.59.1

.PHONY: lint
lint: .install-linter
	### RUN GOALNGCI-LINT ###
	$(GOLANGCI_LINT) run ./... --congif=./.golangci.yml

.PHONY: docker-up
docker-up:
	### START DOCKER COMPOSE ###
	docker compose -f ./docker/docker-compose.yml up -d

.PHONY: docker-down
docker-down:
	### STOP DOCKER COMPOSE ###
	docker compose -f ./docker/docker-compose.yml down

.PHONY: docker-restart
docker-restart: docker-down docker-up

.PHONY: docker-logs
docker-logs:
	### VIEW DOCKER COMPOSE LOGS ###
	docker compose logs -f

install-deps:
	GOBIN=$(PROJECT_BIN) go install "google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2"
	GOBIN=$(PROJECT_BIN) go install "google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1"

get-deps:
	go get -u "google.golang.org/protobuf/cmd/protoc-gen-go"
	go get -u "google.golang.org/grpc/cmd/protoc-gen-go-grpc"

.PHONY: make generate-metrics-api
generate-metrics-api:
	mkdir -p pkg/metrics_v1
	protoc --proto_path proto \
	--go_out=pkg/metrics_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/metrics_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	proto/go-metrics-altering.proto

# Run tests with coverage only for pkg and internal folders
# Folders cmd/server, cmd/agent, internal/agent/app, pkg/metrics_v1 were excluded because of:
# 1. Folder internal/agent/app contains code which is responible for initialazing funcs for sending metrics.
# All such funcs were tested in relevant sections.
# 2. Folders cmd/server and cmd/agent contains code for starting up agent and server accordingy.
# 3. Folder metrics_v1 contains generated proto files.
.PHONY: make test-cover
test-cover:
	go test -cover -v -coverpkg=./pkg/crypt...,./internal/agent/memory...,./internal/agent/sendMetrics...,./internal/server/http...,./internal/server/grpc/rpcserver...,./internal/storage/inmemory...,./internal/storage/postgres...,./pkg/httpServer...,./pkg/interceptors...,./pkg/logger...,./pkg/middlewares...,./pkg/processJSON...,./pkg/processMap... -coverprofile=profile.cov ./...
	go tool cover -func profile.cov
