MIGRATION_FOLDER := sql/migrations/
POSTGRES_URL ?= postgres:///ovh?sslmode=disable

clean:
	go clean -testcache

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

test:
	go test -v ./...

test-single:
	go test -v ./... -testify.m $(name)

updater:
	go run ./cmd/updater
