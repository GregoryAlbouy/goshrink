.PHONY: start
start:
	@make docker && \
	make run

.PHONY: run
run:
	@go run cmd/server/main.go

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

.PHONY: docker
docker:
	@docker-compose --env-file ./.env up --build --detach

.PHONY: docker-down
docker-down:
	@docker-compose --env-file ./.env down

.PHONY: docker-restart
docker-restart:
	@docker-compose --env-file ./.env restart

.PHONY: start-queue
start-queue:
	@docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management

.PHONY: start-server
start-server:
	@go run cmd/server/main.go

.PHONY: start-server-migrate
start-server-migrate:
	@go run cmd/server/main.go -m

.PHONY: start-storage
start-storage:
	@go run $$(ls -1 ./cmd/storage/*.go | grep -v _test.go)

.PHONY: start-worker
start-worker:
	@go run $$(ls -1 ./cmd/worker/*.go | grep -v _test.go)

.PHONY: post-guest
post-guest:
	curl -X POST -H "Content-Type:application/json" -d '{"username": "guest", "email": "guest@goshrink.com", "password": "password"}' http://localhost:9999/users

.PHONY: login
login:
	 curl -X POST -H "Content-Type: application/json" -d '{"username": "Bret", "password": "password"}' http://localhost:9999/login

.PHONY: post-avatar
post-avatar:
	@curl -X POST -H "Authorization:Bearer ${t}" -H "Content-Type:multipart/form-data" -F "image=@fixtures/sample.png" http://localhost:9999/users/1/avatar

.PHONY: test
test:
	@go test -v -timeout 30s -run ${t} ./...

.PHONY: tests
tests:
	@go test -v -timeout 30s ./...

.PHONY: docs
docs:
	@godoc -http=localhost:9995
