BIN := bin/indexer

ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: up down wait-db migrate build run dev clean

up:
	docker-compose up -d

down:
	docker-compose down

wait-db:
	@echo "Waiting for postgres..."
	@until docker exec erc20_indexer_db pg_isready -U postgres > /dev/null 2>&1; do sleep 1; done
	@echo "Postgres ready"

migrate: wait-db
	docker exec -i erc20_indexer_db psql -U postgres -d indexer < migrations/001_init.sql

build:
	go build -o $(BIN) ./cmd/indexer

run: build
	./$(BIN)

dev: up migrate run

clean:
	rm -f $(BIN)
