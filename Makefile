.PHONY: help generate init-schema
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

init-schema: ## Initialize a new schema
ifdef name
	@go run entgo.io/ent/cmd/ent init $$name
else
	@echo "Usage: make init-schema name=SchemaName" && exit 64
endif
