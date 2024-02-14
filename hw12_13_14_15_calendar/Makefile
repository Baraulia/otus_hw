#General envs
DATABASE ?= sql
TRANSPORT ?= http
#Postrges
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= password
POSTGRES_DB ?= backend
POSTGRES_PORT ?= 5435
POSTGRES_CONTAINER := postgres-calendar
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
RABBIT_CONTAINER := rabbitmq

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build-calendar:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar
build-scheduler:
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
build-sender:
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run-postgres:
	docker run -d --name $(POSTGRES_CONTAINER) \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-e POSTGRES_DB=$(POSTGRES_DB) \
	-p $(POSTGRES_PORT):5432 \
	-v postgres-data:/var/lib/postgresql/data \
	postgres:latest

run-rabbit:
	docker run -d --name $(RABBIT_CONTAINER) \
	-p $(RABBIT_PORT):5672 \
	-p $(RABBIT_PORT_UI):15672 \
	-e RABBITMQ_DEFAULT_USER=$(RABBIT_USER) \
	-e RABBITMQ_DEFAULT_PASS=$(RABBIT_PASSWORD) \
	rabbitmq:management

run: build-calendar
ifeq ($(DATABASE),sql)
ifneq ($(shell docker ps -q --filter "name=$(POSTGRES_CONTAINER)"),)
	@echo "Container $(POSTGRES_CONTAINER) is already running."
else
	@$(MAKE) run-postgres
endif
endif

ifneq ($(shell docker ps -q --filter "name=$(RABBIT_CONTAINER)"),)
	@echo "Container $(RABBIT_CONTAINER) is already running."
else
	@$(MAKE) run-rabbit
endif
	$(BIN_CALENDAR) -config ./configs/config.yaml -database $(DATABASE) -transport $(TRANSPORT)

run-scheduler: build-scheduler
	$(BIN_SCHEDULER) -config ./configs/scheduler_config.yaml -database $(DATABASE)
run-sender: build-sender
	$(BIN_SENDER) -config ./configs/sender_config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(CALENDAR_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(CALENDAR_IMG)

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

.PHONY: build run build-img run-img version test lint
