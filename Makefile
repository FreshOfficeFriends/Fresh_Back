include .env

.DEFAULT_GOAL = run

LOCAL_BIN=$(CURDIR)/bin

run:
	go run cmd/sso/main.go

build:
	CGO_ENABLED=0 GOOS=linux go build -o freshFriends cmd/sso/main.go

install_deps:
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

#by default is checked using pre-commit
lint:
	$(LOCAL_BIN)/golangci-lint run --fast --config .golangci.pipeline.yaml

migration-status:
	$(LOCAL_BIN)/migrate -source file://migrations -database $(DB_DSN) version

migration-up:
	$(LOCAL_BIN)/migrate -source file://migrations -database $(DB_DSN) up 1

migration-down:
	$(LOCAL_BIN)/migrate -source file://migrations -database $(DB_DSN) down 1


dockerFull:
	docker-compose -f docker-compose.full_depend.yml build --no-cache && docker-compose -f docker-compose.full_depend.yml up -d

.PHONY: run install_deps lint migration-status migration-up migration-down docker build
