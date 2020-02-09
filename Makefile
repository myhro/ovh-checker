MIGRATION_FOLDER := sql/migrations/
POSTGRES_URL ?= postgres:///ovh?sslmode=disable

create:
	migrate create -dir $(MIGRATION_FOLDER) -ext sql -seq $(name)

destroy:
	migrate -database $(POSTGRES_URL) -path $(MIGRATION_FOLDER) down

migrate:
	migrate -database $(POSTGRES_URL) -path $(MIGRATION_FOLDER) up

updater:
	go run ./cmd/updater
