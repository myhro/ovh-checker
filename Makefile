COVERAGE := coverage.out
COVERAGE_REPORT := coverage.html
MIGRATION_FOLDER := sql/migrations/
POSTGRES_URL ?= postgres:///ovh?sslmode=disable

.PHONY: api

api:
	go run ./api

clean:
	go clean -testcache
	rm -f $(COVERAGE) $(COVERAGE_REPORT)

coverage:
	go tool cover -html=$(COVERAGE) -o $(COVERAGE_REPORT)

create:
	migrate create -dir $(MIGRATION_FOLDER) -ext sql -seq $(name)

destroy:
	migrate -database $(POSTGRES_URL) -path $(MIGRATION_FOLDER) down

lint:
	@golint -set_exit_status ./...

migrate:
	migrate -database $(POSTGRES_URL) -path $(MIGRATION_FOLDER) up

notifier:
	go run ./cmd/notifier

test: test-all coverage

test-all:
	go test -coverprofile $(COVERAGE) -v ./...

test-single:
	go test -v ./... -testify.m $(name)

updater:
	go run ./cmd/updater
