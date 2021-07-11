BIN := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run: build
	$(BIN) -config ./configs/config.json

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race -v ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.37.0

lint: install-lint-deps
	golangci-lint run

generate-deps:
	mkdir -p internal/server/eventpb

generate: generate-deps
	protoc -I ./proto \
    --go_out ./internal/server/eventpb/ --go_opt paths=source_relative \
    --go-grpc_out ./internal/server/eventpb/ --go-grpc_opt paths=source_relative \
    proto/EventService/*.proto

generate-gateway: generate
	protoc -I ./proto --grpc-gateway_out ./internal/server/eventpb \
			--grpc-gateway_opt logtostderr=true \
			--grpc-gateway_opt paths=source_relative \
			--grpc-gateway_opt generate_unbound_methods=true \
			proto/EventService/*.proto

up:
	docker-compose up -d

upl:
	docker-compose up

down:
	docker-compose down

rebuild:
	docker-compose up --build

integration-tests-cleanup:
	docker-compose -f docker-compose.test.yaml down \
    --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \
  	docker-compose rm -f

integration-tests:
	set -e ;\
	docker-compose -f docker-compose.test.yaml up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f docker-compose.test.yaml run integration_tests go test -v || test_status_code=$$? ;\
	docker-compose -f docker-compose.test.yaml down ;\
	exit $$test_status_code ;

.PHONY: build run build-img run-img version test lint
