# DEVELOPMENT

MIGRATE=docker compose run --rm migrate


# Create a new migration
migrate-create:
	@read -p "Enter migration name: " name; \
	$(MIGRATE) create "$$name" sql
	goose -dir .sql/migrations create "$$name" sql

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
	docker compose run --rm migrate up

.PHONY: migrate-down
migrate-down:
	docker compose run --rm migrate down
