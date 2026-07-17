include .env
export


# For UNIX
# export PROJECT_ROOT=$(shell pwd)
# For Windows
export PROJECT_ROOT=C:/Users/Svat/Documents/Github/golang-todo
export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs


env-up:
	docker compose up -d todo-app-postgres

env-down:
	docker compose down todo-app-postgres

env-cleanup:
	@read -p "Are you sure you want to cleanup all environment? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todo-app-postgres port-forwarder && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Cleanup complete"; \
	else \
		echo "Cleanup aborted"; \
	fi

env-cleanup-win:
	@docker compose down todo-app-postgres port-forwarder | \
	rmdir /s /q "${PROJECT_ROOT}/out/pgdata" | \
	echo Cleanup complete

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

migrate-create-win:
	@docker compose run --rm todo-app-postgres-migrate \
		create -dir /migrations -ext sql -seq "$(name)"

migrate-up:
	@make migrate-action action=up

migrate-up-win:
	@docker compose run --rm todo-app-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todo-app-postgres:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable \
		up

migrate-down:
	@make migrate-action action=down

migrate-down-win:
	@docker compose run --rm todo-app-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todo-app-postgres:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable \
		down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Action is required. Use 'make migrate-action action=<action>' to specify an action."; \
		exit 1; \
	fi; \
	docker compose run --rm todo-app-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todo-app-postgres:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

logs-cleanup:
	@read -p "Are you sure you want to cleanup logs? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${LOGGER_FOLDER} && \
		echo "Cleanup complete"; \
	else \
		echo "Cleanup aborted"; \
	fi

logs-cleanup-win:
	@rmdir /s /q "${LOGGER_FOLDER}" | \
	echo Cleanup complete

app-run:
	@export go mod tidy && \
	go run cmd/main.go

app-run-win:
	@go mod tidy | \
	go run cmd/main.go

app-deploy:
	@docker compose up -d --build todo-app
