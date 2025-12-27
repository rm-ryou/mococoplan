#! /usr/bin/env -S just --justfile

set dotenv-load

DB_PORT := "3306"
DB_HOST := env("MYSQL_HOST", "mysql")
DB_NETWORK := env("MYSQL_NETWORK", "mococoplan_default")

DB_NAME := env("MYSQL_DATABASE", "mococoplan")
DB_USER := env("MYSQL_USER", "user")
DB_PASSWORD := env("MYSQL_PASSWORD", "password")

MIGRATIONS_DIR := "db/migrations"
MIGRATE_IMAGE := "migrate/migrate:v4.19.0"

default:
  @just --list

[group: "db"]
migrate-create NAME:
  mkdir -p {{MIGRATIONS_DIR}}
  docker run --rm \
    -v "{{justfile_directory()}}/{{MIGRATIONS_DIR}}:/migrations" \
    {{MIGRATE_IMAGE}} \
    create -ext sql -dir /migrations -seq {{NAME}}

[group: "db"]
migrate-up:
  docker run --rm \
    -v {{justfile_directory()}}/{{MIGRATIONS_DIR}}:/migrations \
    --network {{DB_NETWORK}} \
    {{MIGRATE_IMAGE}} \
    -path=/migrations \
    -database "mysql://{{DB_USER}}:{{DB_PASSWORD}}@tcp({{DB_HOST}}:{{DB_PORT}})/{{DB_NAME}}" \
    up

[group: "db"]
migrate-down:
  docker run --rm \
    -v {{justfile_directory()}}/{{MIGRATIONS_DIR}}:/migrations \
    --network {{DB_NETWORK}} \
    {{MIGRATE_IMAGE}} \
    -path=/migrations \
    -database "mysql://{{DB_USER}}:{{DB_PASSWORD}}@tcp({{DB_HOST}}:{{DB_PORT}})/{{DB_NAME}}" \
    down -all

[group: "db"]
migrate-down-1:
  docker run --rm \
    -v {{justfile_directory()}}/{{MIGRATIONS_DIR}}:/migrations \
    --network {{DB_NETWORK}} \
    {{MIGRATE_IMAGE}} \
    -path=/migrations \
    -database "mysql://{{DB_USER}}:{{DB_PASSWORD}}@tcp({{DB_HOST}}:{{DB_PORT}})/{{DB_NAME}}" \
    down 1

[group: "db"]
migrate-force VERSION:
  docker run --rm \
    -v {{justfile_directory()}}/{{MIGRATIONS_DIR}}:/migrations \
    --network {{DB_NETWORK}} \
    {{MIGRATE_IMAGE}} \
    -path=/migrations \
    -database "mysql://{{DB_USER}}:{{DB_PASSWORD}}@tcp({{DB_HOST}}:{{DB_PORT}})/{{DB_NAME}}" \
    force {{VERSION}}

[group: "db"]
migrate-version:
  # TODO: avoid hard cording
  docker run --rm \
    -v {{justfile_directory()}}/{{MIGRATIONS_DIR}}:/migrations \
    --network {{DB_NETWORK}} \
    {{MIGRATE_IMAGE}} \
    -path=/migrations \
    -database "mysql://{{DB_USER}}:{{DB_PASSWORD}}@tcp({{DB_HOST}}:{{DB_PORT}})/{{DB_NAME}}" \
    version

[group: "db"]
exec-db:
  docker compose exec mysql mysql -u {{DB_USER}} -p mococoplan

[group: "server"]
up:
  docker compose up --build --detach

[group: "server"]
down:
  docker compose down

[group: "server"]
logs:
  docker compose logs -f

help:
  @echo "Availabel recipes:"
  @echo ""
  @echo "Server"
  @echo "  just up                        # Start Docker containers"
  @echo "  just down                      # Stop Docker containers"
  @echo "  just logs                      # Show Docker containers"
  @echo ""
  @echo "DB"
  @echo "  just migrate-create <name>     # Create new migration file"
  @echo "  just migrate-up                # Run all pending migrations"
  @echo "  just migrate-down              # Rollback all migrations"
  @echo "  just migrate-down-1            # Rollback one migration"
  @echo "  just migrate-force <version>   # Force set migration version"
  @echo "  just migrate-version           # Show current migrate"
