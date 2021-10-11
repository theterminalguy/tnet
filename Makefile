.PHONY: help generate scaffold
.DEFAULT_GOAL: help

default: help

help: ## Output available commands
	@echo "Available commands:"
	@echo
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

generate: ## Generate ent Assests
	@go generate ./ent

start: ## Start the app
	@go run .

scaffold: ## Generate a new resource scaffold
ifdef resource
	@go run cmd/generate/main.go $$resource
else
	@echo "Usage: make scaffold resource=resource_name" && exit 64
endif
