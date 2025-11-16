#! /usr/bin/env -S just --justfile

MIGRATIONS_DIR := "db/migrations"
MIGRATE_IMAGE := "migrate/migrate:v4.19.0"

default:
  @just --list

help:
  @echo "Availabel recipes:"
  @echo "  just migrate-create name=<name>  # create new migration file"

[group: 'db']
migrate-create NAME:
  mkdir -p {{MIGRATIONS_DIR}}
  docker run --rm \
    -v "{{justfile_directory()}}/{{MIGRATIONS_DIR}}:/migrations" \
    {{MIGRATE_IMAGE}} \
    create -ext sql -dir /migrations -seq {{NAME}}
