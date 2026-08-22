# ------------------------------------------------------------------
# TickStorm - automação de desenvolvimento
# ------------------------------------------------------------------
SHELL := /bin/sh
BACKEND_DIR := backend
MIGRATIONS_DIR := $(BACKEND_DIR)/database/migrations

# Carrega o .env (se existir) para que DATABASE_URL e afins fiquem disponíveis.
ifneq (,$(wildcard ./.env))
include .env
export
endif

POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5432
POSTGRES_USER ?= root
POSTGRES_PASSWORD ?= toor
POSTGRES_DB ?= tickstorm
POSTGRES_SSLMODE ?= disable

# Derived from the parts above so .env stays the single source of truth.
DATABASE_URL ?= postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)

.DEFAULT_GOAL := help
.PHONY: help tools env infra-up infra-down infra-logs infra-reset infra-ps \
        run build tidy db-gen migrate-up migrate-down migrate-create migrate-force psql redis-cli

help: ## Lista os comandos disponíveis
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

# ---------------------------- Setup -------------------------------
tools: ## Instala as ferramentas de desenvolvimento (air, sqlc, migrate)
	go install github.com/air-verse/air@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

env: ## Cria o .env a partir do .env.example (não sobrepõe um existente)
	@test -f .env && echo ".env já existe, nada a fazer." || (cp .env.example .env && echo ".env criado.")

# -------------------------- Infraestrutura ------------------------
infra-up: ## Arranca Postgres, Redis, Prometheus e Grafana
	docker compose up -d

infra-down: ## Desliga a infraestrutura (mantém os volumes)
	docker compose down

infra-reset: ## Desliga e APAGA todos os dados (volumes incluídos)
	docker compose down -v

infra-ps: ## Estado dos containers
	docker compose ps

infra-logs: ## Segue os logs da infraestrutura
	docker compose logs -f

# --------------------------- Aplicação ----------------------------
run: ## Arranca o backend com live-reload (Air)
	cd $(BACKEND_DIR) && air

build: ## Compila o binário da API
	cd $(BACKEND_DIR) && go build -o ../bin/api ./cmd/api

tidy: ## Sincroniza as dependências do módulo
	cd $(BACKEND_DIR) && go mod tidy

# --------------------------- Base de dados ------------------------
db-gen: ## Gera o código type-safe a partir do SQL (sqlc)
	cd $(BACKEND_DIR) && sqlc generate

migrate-up: ## Aplica todas as migrações pendentes
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down: ## Reverte a última migração
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

migrate-force: ## Desbloqueia um estado "dirty" (uso: make migrate-force version=N)
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(version)

migrate-create: ## Cria um novo par de migrações (uso: make migrate-create name=add_x)
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -format 20060102150405 $(name)

# ----------------------------- Atalhos ----------------------------
psql: ## Abre um shell psql no container do Postgres
	docker compose exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

redis-cli: ## Abre um shell redis-cli no container do Redis
	docker compose exec redis redis-cli
