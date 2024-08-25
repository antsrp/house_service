ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: --setpath, migrateup, migratedown, update, run, docker-start, docker-stop, docker-up, docker-down, run-all, build-run, run-builded, build

update:
	go get ./...
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

run-all: docker-up update migrateup build-run

run:
	@go run ./cmd/api

build-run: update build run-builded
run-builded: 
	./bin/main
build:
	@go build -o bin/main ./cmd/api

test:
	go test ./...

migrateup: --setpath
	migrate -path db/migrations -database $(dbpath) -verbose up

migratedown: --setpath
	migrate -path db/migrations -database $(dbpath) -verbose down

docker-up:
	docker compose build --no-cache
	docker compose up -d

docker-down:
	docker compose down

docker-start:
	docker compose start

docker-stop:
	docker compose stop

--setpath:
	$(eval dbpath = $(DB_TYPE)://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable)