default: help

help: ## Output available commands
	@echo "Available commands:"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

setup: ## Builds the web container
	STAGE=app-build docker-compose -f docker-compose.yml build web

start: ## Start all services
	STAGE=app-build docker-compose -f docker-compose.yml up

stop: ## Stop all services
	STAGE=app-build docker-compose -f docker-compose.yml down

destroy: ## Remove all containers and images. Also, destroy all volumes
	STAGE=app-build docker-compose -f docker-compose.yml down -v --remove-orphans --rmi all

scaffold: ## Generate a new resource scaffold
	STAGE=app-build docker-compose run web go run cmd/generate/main.go $(resource)

test: ## Run all tests
	STAGE=tests docker-compose -f docker-compose.yml build web 
	STAGE=tests docker-compose -f docker-compose.yml up web

hot-reload: ## Enables hot reload for the web service
	STAGE=hot-reload docker-compose -f docker-compose.yml build web 
	STAGE=hot-reload docker-compose -f docker-compose.yml up web
