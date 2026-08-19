include .env
export


export PROJECT_ROOT=$(shell pwd)
export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs


env-up:
	docker compose up -d cascade-app-postgres cascade-app-redis cascade-app-mongo

env-down:
	docker compose down cascade-app-postgres cascade-app-redis cascade-app-mongo

env-cleanup:
	@read -p "Are you sure you want to cleanup all environment? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down cascade-app-postgres cascade-app-redis cascade-app-mongo port-forwarder && \
		sudo rm -rf ${PROJECT_ROOT}/out/pg_data && \
		sudo rm -rf ${PROJECT_ROOT}/out/redis_data && \
		sudo rm -rf ${PROJECT_ROOT}/out/mongo_data && \
		echo "Cleanup complete"; \
	else \
		echo "Cleanup aborted"; \
	fi

env-cleanup-win:
	@docker compose down cascade-app-postgres cascade-app-redis cascade-app-mongo port-forwarder | \
	rmdir /s /q "${PROJECT_ROOT}/out/pg_data" | \
	rmdir /s /q "${PROJECT_ROOT}/out/redis_data" | \
	rmdir /s /q "${PROJECT_ROOT}/out/mongo_data" | \
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
	docker compose run --rm cascade-app-postgres-migrate \
		create -dir /migrations -ext sql -seq "$(name)"

migrate-create-win:
	@docker compose run --rm cascade-app-postgres-migrate \
		create -dir /migrations -ext sql -seq "$(name)"

migrate-up:
	@make migrate-action action=up

migrate-up-win:
	@docker compose run --rm cascade-app-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@cascade-app-postgres:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable \
		up

migrate-down:
	@make migrate-action action=down

migrate-down-win:
	@docker compose run --rm cascade-app-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@cascade-app-postgres:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable \
		down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Action is required. Use 'make migrate-action action=<action>' to specify an action."; \
		exit 1; \
	fi; \
	docker compose run --rm cascade-app-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@cascade-app-postgres:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

mongo-init:
	@docker exec -i cascade-app-mongo-database mongosh \
		-u ${MONGO_USER} \
		-p ${MONGO_PASSWORD} \
		--authenticationDatabase admin \
		< scripts/mongo-init.js

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

swagger-gen:
	@docker compose run --rm swagger \
		init \
		-g cmd/main.go \
		-o docs \
		--parseInternal \
		--parseDependency

app-run:
	@export go mod tidy && \
	go run cmd/main.go

app-run-win:
	@go mod tidy | \
	go run cmd/main.go

app-deploy:
	@docker compose up -d --build cascade-app
