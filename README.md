# Link Shortener

## Prerequisites
- Docker & Docker compose

## Quick Start
docker-compose up

## Services
- App: localhost:(localport)/main/
- DB: postgres:5432

## Env Variables
1. Create `.env` with:
PORT=...
GOOSE_DBSTRING=...
DB_URL=...
GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=./sql/migrations
2. Create `postgres.env` with:
POSTGRES_PASSWORD=
POSTGRES_USER=
POSTGRES_DB=