.PHONY: up down build logs nuke

up:
	@echo "🚀 Starting LeetFalls (with build) ..."
	docker-compose up --build

down:
	@echo "🛑 Stopping and removing containers ..."
	docker-compose down

logs:
	@echo "📜 Streaming logs (latest 50 lines) ..."
	docker-compose logs -f --tail=50

nuke:
	@echo "💣 Nuking and rebuilding everything from scratch ..."
	docker-compose down -v --remove-orphans
	docker-compose up --build
