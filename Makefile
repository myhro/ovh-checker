MIGRATION_FOLDER := sql/migrations/
POSTGRES_URL ?= postgres:///ovh?sslmode=disable

clean:
	go clean -testcache

create:
	migrate create -dir $(MIGRATION_FOLDER) -ext sql -seq $(name)

destroy:
	migrate -database $(POSTGRES_URL) -path $(MIGRATION_FOLDER) down

migrate:
	migrate -database $(POSTGRES_URL) -path $(MIGRATION_FOLDER) up

test:
	go test -v ./...

updater:
	go run ./cmd/updater
