# Variáveis
DB_URL="postgresql://root:toor@localhost:5432/tickstorm?sslmode=disable"

# Iniciar infraestrutura (BD, Redis, Grafana, Prometheus)
infra-up:
	docker-compose up -d

# Desligar infraestrutura
infra-down:
	docker-compose down

# Arrancar o backend com live-reload (Air)
run:
	air

# Gerar código do banco de dados (sqlc)
db-gen:
	sqlc generate

# Instalar as ferramentas essenciais de desenvolvimento
tools:
	go install github.com/cosmtrek/air@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest