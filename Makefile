.PHONY: default
default:
	@make docker-up

# Full app build in docker

.PHONY: docker-up
docker-up:
	@docker-compose --env-file ./.env up --build

.PHONY: docker-down
docker-down:
	@docker-compose --env-file ./.env down

.PHONY: docker-restart
docker-restart:
	@docker-compose --env-file ./.env restart

# Individual docker builds

.PHONY: server
server:
	@docker-compose --env-file ./.env up --build server

.PHONY: mysql
mysql:
	@docker-compose --env-file ./.env up --build mysql

.PHONY: storage
storage:
	@docker-compose --env-file ./.env up --build storage

.PHONY: worker
worker:
	@docker-compose --env-file ./.env up --build worker

.PHONY: rabbitmq
rabbitmq:
	@docker-compose --env-file ./.env up --build rabbitmq

# Individual local builds

.PHONY: local-server
local-server:
	@go run $$(ls -1 ./cmd/server/*.go | grep -v _test.go)

.PHONY: local-storage
local-storage:
	@go run $$(ls -1 ./cmd/storage/*.go | grep -v _test.go)

.PHONY: local-worker
local-worker:
	@go run $$(ls -1 ./cmd/local/*.go | grep -v _test.go)

# Test commands

.PHONY: test
test:
	@go test -v -timeout 30s -run ${t} ./...

.PHONY: tests
tests:
	@go test -v -timeout 30s ./...

# Serve docs

.PHONY: docs
docs:
	@godoc -http=localhost:9995
