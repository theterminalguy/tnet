default: help

help: ## Output available commands
	@echo "Available commands:"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

setup: ## Builds the web container
	STAGE=app-build docker-compose -f docker-compose.yml build web

start: ## Start all services
	STAGE=app-build docker-compose -f docker-compose.yml up

dev: ## Run the web server in dev mode without using docker
	ENV=dev go run cmd/web/main.go

stop: ## Stop all services
	STAGE=app-build docker-compose -f docker-compose.yml down

destroy: ## Remove all containers and images. Also, destroy all volumes
	STAGE=app-build docker-compose -f docker-compose.yml down -v --remove-orphans --rmi all

scaffold: ## Generate a new resource scaffold
	go run cmd/generate/main.go $(resource)

test: ## Run all tests
	STAGE=tests docker-compose -f docker-compose.yml build web 
	STAGE=tests docker-compose -f docker-compose.yml up web

hot-reload: ## Enables hot reload for the web service
	STAGE=hot-reload docker-compose -f docker-compose.yml build web 
	STAGE=hot-reload docker-compose -f docker-compose.yml up web

ent-generate: ## Generate ent Assests
	go generate ./ent

go-format: ## Run go fmt ./... on all go files
	go fmt ./...
