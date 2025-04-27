# Makefile for Insightly (new architecture)

COMPOSE = docker-compose

# Start Redis and Postgres only
up:
	@echo "🔵 Starting Redis and Postgres..."
	$(COMPOSE) up -d

# Stop Redis and Postgres and kill running apps
down:
	@echo "🔴 Stopping Redis, Postgres and local app binaries..."
	$(COMPOSE) down
	@pkill -f ./bin/ingest || echo "🛡️ Ingest not running."
	@pkill -f ./bin/summarizer || echo "🛡️ Summarizer not running."
	@pkill -f ./bin/query || echo "🛡️ Query not running."

# View Redis and Postgres logs
logs:
	@echo "🪵 Showing logs from Redis and Postgres..."
	docker logs -f insightly_redis &
	docker logs -f insightly_postgres

# Build and run Ingest Server
ingest:
	@echo "🚀 Building and starting Ingest Server..."
	go build -o ./bin/ingest cmd/ingest/main.go
	./bin/ingest

# Build and run Summarizer Worker
summarizer:
	@echo "⚡ Building and starting Summarizer Worker..."
	go build -o ./bin/summarizer cmd/summarizer/main.go
	./bin/summarizer

# Build and run Query Server
query:
	@echo "🔍 Building and starting Query Server..."
	go build -o ./bin/query cmd/query/main.go
	./bin/query

# Clean Docker system
prune:
	@echo "⚠️ Cleaning Docker system..."
	docker system prune -af --volumes