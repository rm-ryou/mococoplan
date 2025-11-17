#! /usr/bin/env -S just --justfile

set dotenv-load

DB_HOST := "localhost"
DB_PORT := "3306"
DB_NAME := env("MYSQL_DATABASE", "mococoplan")
DB_USER := env("MYSQL_USER", "user")
DB_PASSWORD := env("MYSQL_PASSWORD", "password")

MIGRATIONS_DIR := "db/migrations"
MIGRATE_IMAGE := "migrate/migrate:v4.19.0"

default:
  @just --list

help:
  @echo "Availabel recipes:"
  @echo "  just migrate-create name=<name>  # create new migration file"

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
    --network "mococoplan_default" \
    {{MIGRATE_IMAGE}} \
    -path=/migrations \
    -database "mysql://{{DB_USER}}:{{DB_PASSWORD}}@tcp(mysql:{{DB_PORT}})/{{DB_NAME}}" \
    up

[group: "db"]
migrate-down:
  docker run --rm \
    -v {{justfile_directory()}}/{{MIGRATIONS_DIR}}:/migrations \
    --network "mococoplan_default" \
    {{MIGRATE_IMAGE}} \
    -path=/migrations \
    -database "mysql://{{DB_USER}}:{{DB_PASSWORD}}@tcp(mysql:{{DB_PORT}})/{{DB_NAME}}" \
    down -all

[group: "db"]
migrate-down-1:
  docker run --rm \
    -v {{justfile_directory()}}/{{MIGRATIONS_DIR}}:/migrations \
    --network "mococoplan_default" \
    {{MIGRATE_IMAGE}} \
    -path=/migrations \
    -database "mysql://{{DB_USER}}:{{DB_PASSWORD}}@tcp(mysql:{{DB_PORT}})/{{DB_NAME}}" \
    down 1

[group: "db"]
migrate-force VERSION:
  docker run --rm \
    -v {{justfile_directory()}}/{{MIGRATIONS_DIR}}:/migrations \
    --network "mococoplan_default" \
    {{MIGRATE_IMAGE}} \
    -path=/migrations \
    -database "mysql://{{DB_USER}}:{{DB_PASSWORD}}@tcp(mysql:{{DB_PORT}})/{{DB_NAME}}" \
    force {{VERSION}}

[group: "db"]
migrate-version:
  # TODO: avoid hard cording
  docker run --rm \
    -v {{justfile_directory()}}/{{MIGRATIONS_DIR}}:/migrations \
    --network "mococoplan_default" \
    {{MIGRATE_IMAGE}} \
    -path=/migrations \
    -database "mysql://{{DB_USER}}:{{DB_PASSWORD}}@tcp(mysql:{{DB_PORT}})/{{DB_NAME}}" \
    version
