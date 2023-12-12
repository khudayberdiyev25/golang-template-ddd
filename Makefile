-include .env
export

CURRENT_DIR=$(shell pwd)
APP=iman_telegram_service
CMD_DIR=./cmd
DOCKER_COMPOSE_FILE=docker-compose.yaml
.DEFAULT_GOAL = build

arg = $(filter-out $@,$(MAKECMDGOALS))

# go generate
.PHONY: go-gen
go-gen:
	go generate ./...

# build for current os
.PHONY: build
build:
	go build -ldflags="-s -w" -o ./bin/${APP} ${CMD_DIR}/main.go

# build for linux amd64
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOARCH="amd64" GOOS=linux go build -ldflags="-s -w" -o ./bin/${APP} ${CMD_DIR}/main.go

.PHONY: start
start:
	@echo "Start Containers"
	docker-compose -f ${DOCKER_COMPOSE_FILE} up -d ${DOCKER_SERVICES}
	sleep 2
	docker-compose -f ${DOCKER_COMPOSE_FILE} ps

.PHONY: stop
stop:
	@echo "Stop Containers"
	docker-compose -f ${DOCKER_COMPOSE_FILE} stop ${DOCKER_SERVICES}
	sleep 2
	docker-compose -f ${DOCKER_COMPOSE_FILE} ps

.PHONY: stop
rm: stop
	@echo "Remove Containers"
	docker-compose -f ${DOCKER_COMPOSE_FILE} rm -v -f ${DOCKER_SERVICES}

.PHONY: migration-up
migration-up:
	@echo "Migrations Up"
	sleep 2
	docker-compose run --rm migrate -path=migrations/ -database='postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable' up

.PHONY: migration-generate
migration-generate:
	@echo "Generation migration file $(name)"
	sleep 2
	docker-compose run --rm migrate create -ext sql -dir ./migrations -seq $(name)

.PHONY: mod-download
mod-download:
	@echo "Go mod download"
	sleep 2
	docker-compose exec grpc go mod download

.PHONY: go-get
go-get:
	@echo "Go get ${arg}"
	sleep 2
	docker-compose exec grpc go get -d ${arg}
