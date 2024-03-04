#General envs
DATABASE ?= sql
TRANSPORT ?= http
#Postrges
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= password
POSTGRES_DB ?= backend
POSTGRES_PORT ?= 5435
POSTGRES_CONTAINER ?= postgres-calendar
#Golang
BIN_CALENDAR := "./bin/calendar"
BIN_SCHEDULER := "./bin/scheduler"
BIN_SENDER := "./bin/sender"
CALENDAR_IMG="calendar:develop"
SCHEDULER_IMG="scheduler:develop"
SENDER_IMG="sender:develop"
#Rabbit
RABBIT_USER ?= rabbit
RABBIT_PASSWORD ?= password
RABBIT_PORT ?= 5672
RABBIT_PORT_UI ?= 15672
RABBIT_CONTAINER ?= rabbitmq

CALENDAR_CONTAINER ?= calendar
SCHEDULER_CONTAINER ?= calendar-scheduler
SENDER_CONTAINER ?= calendar-sender

NETWORK_NAME ?= calendar_network

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)


build-calendar:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar
build-scheduler:
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
build-sender:
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run-postgres:
ifeq ($(DATABASE),sql)
ifneq ($(shell docker ps -q --filter "name=$(POSTGRES_CONTAINER)"),)
	@echo "Container $(POSTGRES_CONTAINER) is already running."
else
	docker run -d --network=$(NETWORK_NAME) --name $(POSTGRES_CONTAINER) \
    	-e POSTGRES_USER=$(POSTGRES_USER) \
    	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
    	-e POSTGRES_DB=$(POSTGRES_DB) \
    	-p $(POSTGRES_PORT):5432 \
    	-v postgres-data:/var/lib/postgresql/data \
    	postgres:latest
endif
endif

run-rabbit:
ifneq ($(shell docker ps -q --filter "name=$(RABBIT_CONTAINER)"),)
	@echo "Container $(RABBIT_CONTAINER) is already running."
else
	docker run -d --network=$(NETWORK_NAME) --name $(RABBIT_CONTAINER) \
	-p $(RABBIT_PORT):5672 \
	-p $(RABBIT_PORT_UI):15672 \
	-e RABBITMQ_DEFAULT_USER=$(RABBIT_USER) \
	-e RABBITMQ_DEFAULT_PASS=$(RABBIT_PASSWORD) \
	rabbitmq:management
endif

run: build-calendar run-postgres run-rabbit
	$(BIN_CALENDAR) -config ./configs/config.yaml -database $(DATABASE) -transport $(TRANSPORT)

run-scheduler: build-scheduler
	$(BIN_SCHEDULER) -config ./configs/scheduler_config.yaml -database $(DATABASE)
run-sender: build-sender
	$(BIN_SENDER) -config ./configs/sender_config.yaml

build-calendar-img:
	docker build \
		--build-arg LDFLAGS="$(LDFLAGS)" \
		-t $(CALENDAR_IMG) \
		-f build/calendar/Dockerfile .

run-calendar-img: build-calendar-img
	docker run --network=$(NETWORK_NAME) --name $(CALENDAR_CONTAINER) \
 	-e DATABASE=$(DATABASE) -e TRANSPORT=$(TRANSPORT) \
 	-e sqlHost=$(POSTGRES_CONTAINER) -e sqlPort=5432 $(CALENDAR_IMG)

build-scheduler-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(SCHEDULER_IMG) \
		-f build/scheduler/Dockerfile .

run-scheduler-img: build-scheduler-img
	docker run --network=$(NETWORK_NAME) \
	--name $(SCHEDULER_CONTAINER) \
	-e sqlHost=$(POSTGRES_CONTAINER) -e sqlPort=5432 \
	-e mbHost=$(RABBIT_CONTAINER) -e mbPort=$(RABBIT_PORT) \
	$(SCHEDULER_IMG)

build-sender-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(SENDER_IMG) \
		-f build/sender/Dockerfile .

run-sender-img:build-sender-img
	docker run --network=$(NETWORK_NAME) \
	--name $(SENDER_CONTAINER) \
	-e mbHost=$(RABBIT_CONTAINER) -e mbPort=$(RABBIT_PORT) \
	 $(SENDER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin latest

lint: install-lint-deps
	golangci-lint run ./...

generate:
	rm -rf internal/api/grpc/pb
	mkdir -p internal/api/grpc/pb
	protoc --proto_path=api/ --go_out=internal/api/grpc/pb	--go-grpc_out=internal/api/grpc/pb api/*.proto

create_network:
	@docker network inspect $(NETWORK_NAME) &>/dev/null || \
    docker network create $(NETWORK_NAME)

delete_network:
	docker network rm $(NETWORK_NAME)

up:
	docker-compose  -f ./deployments/docker-compose.yaml up
down:
	docker-compose  -f ./deployments/docker-compose.yaml down


integration-tests:
	set -e ;\
	docker-compose -f ./deployments/docker-compose.test.yaml up --build -d ;\
	test_status_code=0 ;\
   	docker-compose -f ./deployments/docker-compose.test.yaml run integration_tests go test || test_status_code=$$? ;\
	docker-compose -f ./deployments/docker-compose.test.yaml down ;\
	exit $$test_status_code ;

integration-tests-teardown:
	docker-compose -f ./deployments/docker-compose.test.yaml down \
        --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \
  	docker-compose rm -f

.PHONY: build run build-img run-img version test lint
