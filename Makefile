# DEVELOPMENT

MIGRATE=docker compose run --rm


# Create a new migration
migrate-create:
	@read -p "Enter migration name: " name; \
	goose create "$$name" sql

.PHONY: run
run:
	go run main.go

tidy::
	go tidy

.PHONY: docker-build
docker-build:
	docker compose build

.PHONY: docker-up
docker-up:
	docker compose up -d

.PHONY: docker-db-shell
docker-db-shell:
	docker compose exec -it auth52_db psql -U auth52 -d auth52_db


.PHONY: migrate-up
migrate-up:
	$(MIGRATE) -e GOOSE_COMMAND=up migrate

.PHONY: migrate-down
migrate-down:
	$(MIGRATE) -e GOOSE_COMMAND=down migrate

.PHONY: migrate-down-to
migrate-down-to:
	@read -p "Enter version to roll back to: " version; \
	docker compose run --rm -e GOOSE_COMMAND=down-to -e GOOSE_COMMAND_ARG=$$version migrate


.PHONY: migrate-status
migrate-status:
	$(MIGRATE) -e GOOSE_COMMAND=status migrate