include .env
export


# export PROJECT_ROOT=$(shell pwd)
export PROJECT_ROOT=C:/Users/Svat/Documents/Github/golang-todo


env-up:
	docker compose up -d todo-app-postgres

env-down:
	docker compose down todo-app-postgres

env-cleanup:
	@read -p "Are you sure you want to cleanup? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todo-app-postgres port-forwarder && \
		rm -rf out/pgdata && \
		echo "Cleanup complete"; \
	else \
		echo "Cleanup aborted"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Migration name is required. Use 'make migrate-create name=<name>' to specify a name."; \
		exit 1; \
	fi; \
	docker compose run --rm todo-app-postgres-migrate \
		create -dir /migrations -ext sql -seq "$(name)"

migrate-up:
	@docker compose run --rm todo-app-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todo-app-postgres:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable \
		up
# 	make migrate-action action=up

migrate-down:
	make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Action is required. Use 'make migrate-action action=<action>' to specify an action."; \
		exit 1; \
	fi; \
	docker compose run --rm todo-app-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todo-app-postgres:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

app-run-win:
	@go run cmd/main.go

app-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	go mod tidy && \
	go run cmd/main.go
