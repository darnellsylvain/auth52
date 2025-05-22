# DEVELOPMENT

MIGRATE=docker compose run --rm migrate goose


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
	$(MIGRATE) up

.PHONY: migrate-down
migrate-down:
	$(MIGRATE) down
