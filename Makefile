.PHONY: start
start:
	@make docker && \
	make run

.PHONY: run
run:
	@go run cmd/server/*.go

.PHONY: docker
docker:
	@docker-compose --env-file ./.env up -d

.PHONY: docker-down
docker-down:
	@docker-compose --env-file ./.env down

.PHONY: docker-restart
docker-restart:
	@docker-compose --env-file ./.env restart

.PHONY: test
test:
	@go test -v -timeout 30s -run ${t} ./...

.PHONY: tests
tests:
	@go test -timeout 30s ./...

.PHONY: docs
docs:
	@godoc -http=localhost:9995
