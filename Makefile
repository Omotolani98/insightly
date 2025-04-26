COMPOSE = docker-compose

up:
	@echo "🔵 Starting Insightly services..."
	$(COMPOSE) up --build -d

down:
	@echo "🔴 Stopping Insightly services..."
	$(COMPOSE) down

logs:
	@echo "🪵 Tailing logs from all services..."
	$(COMPOSE) logs -f --tail=100

restart-%:
	@echo "🔁 Restarting service $*..."
	$(COMPOSE) up -d --no-deps --build $*

rebuild:
	@echo "♻️ Rebuilding all services..."
	$(COMPOSE) up -d --build

ps:
	@$(COMPOSE) ps

prune:
	@echo "⚠️ Removing all stopped containers, volumes, networks..."
	docker system prune -af --volumes