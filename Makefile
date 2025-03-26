# DEVELOPMENT
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
	docker exec -it auth52_db psql -U auth52 -d auth52_db

.PHONY: migrate-up
migrate-up:
	docker compose run --rm migrate up
