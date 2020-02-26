API := ./api
API_BINARY := dist/api
BUILD_FLAGS := -ldflags "-s -w"
COVERAGE := coverage.out
COVERAGE_REPORT := coverage.html
DEPLOY_FILE = deploy-k8s.yaml
IMAGE := myhro/ovh-checker
MIGRATION_FOLDER := sql/migrations/
NOTIFIER := ./cmd/notifier
NOTIFIER_BINARY := dist/notifier
POSTGRES_URL ?= postgres:///ovh?sslmode=disable
SESSION_CLEANER := ./cmd/session-cleaner
SESSION_CLEANER_BINARY := dist/session-cleaner
UPDATER := ./cmd/updater
UPDATER_BINARY := dist/updater
VERSION ?= $(shell git rev-parse --short HEAD)

export GOBIN := $(PWD)/.bin

.PHONY: api

api:
	go run $(API)

build: build-api build-notifier build-session-cleaner build-updater

build-api:
	go build $(BUILD_FLAGS) -o $(API_BINARY) $(API)

build-notifier:
	go build $(BUILD_FLAGS) -o $(NOTIFIER_BINARY) $(NOTIFIER)

build-session-cleaner:
	go build $(BUILD_FLAGS) -o $(SESSION_CLEANER_BINARY) $(SESSION_CLEANER)

build-updater:
	go build $(BUILD_FLAGS) -o $(UPDATER_BINARY) $(UPDATER)

clean:
	go clean -testcache
	rm -rf dist/ $(COVERAGE) $(COVERAGE_REPORT) $(DEPLOY_FILE)

coverage:
	go tool cover -html=$(COVERAGE) -o $(COVERAGE_REPORT)

create:
	@$(GOBIN)/migrate create -dir $(MIGRATION_FOLDER) -ext sql -seq $(name)

deploy:
	find k8s/ -name '*.yaml' -exec sed 's/<ENV>/$(ENV)/;s/<HOST>/$(HOST)/;s/<VERSION>/$(VERSION)/' {} \; > $(DEPLOY_FILE)
	kubectl apply -f $(DEPLOY_FILE)

deps:
	go install golang.org/x/lint/golint
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate

destroy:
	@$(GOBIN)/migrate -database $(POSTGRES_URL) -path $(MIGRATION_FOLDER) down

docker:
	docker build -t $(IMAGE) .

lint:
	@$(GOBIN)/golint -set_exit_status ./...

migrate:
	@$(GOBIN)/migrate -database $(POSTGRES_URL) -path $(MIGRATION_FOLDER) up

notifier:
	go run $(NOTIFIER)

push:
	docker tag $(IMAGE):latest $(IMAGE):$(VERSION)
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest

session-cleaner:
	go run $(SESSION_CLEANER)

test: test-all coverage

test-all:
	go test -coverprofile $(COVERAGE) -v ./...

test-single:
	go test $(folder) -testify.m $(name)

test-suite:
	go test -v ./... -run $(name)

updater:
	go run $(UPDATER)
