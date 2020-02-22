API := ./api
API_BINARY := dist/api
COVERAGE := coverage.out
COVERAGE_REPORT := coverage.html
MIGRATION_FOLDER := sql/migrations/
NOTIFIER := ./cmd/notifier
NOTIFIER_BINARY := dist/notifier
POSTGRES_URL ?= postgres:///ovh?sslmode=disable
UPDATER := ./cmd/updater
UPDATER_BINARY := dist/updater

export GOBIN := $(PWD)/.bin

.PHONY: api

api:
	go run $(API)

build: build-api build-notifier build-updater

build-api:
	go build -o $(API_BINARY) $(API)

build-notifier:
	go build -o $(NOTIFIER_BINARY) $(NOTIFIER)

build-updater:
	go build -o $(UPDATER_BINARY) $(UPDATER)

clean:
	go clean -testcache
	rm -rf dist/ $(COVERAGE) $(COVERAGE_REPORT)

coverage:
	go tool cover -html=$(COVERAGE) -o $(COVERAGE_REPORT)

create:
	@$(GOBIN)/migrate create -dir $(MIGRATION_FOLDER) -ext sql -seq $(name)

deps:
	go install golang.org/x/lint/golint
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate

destroy:
	@$(GOBIN)/migrate -database $(POSTGRES_URL) -path $(MIGRATION_FOLDER) down

lint:
	@$(GOBIN)/golint -set_exit_status ./...

migrate:
	@$(GOBIN)/migrate -database $(POSTGRES_URL) -path $(MIGRATION_FOLDER) up

notifier:
	go run $(NOTIFIER)

test: test-all coverage

test-all:
	go test -coverprofile $(COVERAGE) -v ./...

test-single:
	go test -v ./... -testify.m $(name)

updater:
	go run $(UPDATER)
