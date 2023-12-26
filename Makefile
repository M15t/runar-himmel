.DEFAULT_GOAL := help

# Default true if not set
SWAGGER ?= true

# Generates a help message. Borrowed from https://github.com/pydanny/cookiecutter-djangopackage.
help: ## Display this help message
	@echo "Please use \`make <target>' where <target> is one of"
	@perl -nle'print $& if m{^[\.a-zA-Z_-]+:.*?## .*$$}' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-25s\033[0m %s\n", $$1, $$2}'

depends: ## Install & build dependencies
	go mod download

provision: depends ## Provision dev environment
	docker-compose up -d
	sh scripts/waitdb.sh
	@$(MAKE) migrate specs

start: ## Bring up the server on dev environment
	docker-compose up -d
	sh scripts/waitdb.sh
	air

remove.docker: ## Bring down the server on dev environment, remove all docker related stuffs as well
	docker-compose down -v --remove-orphans

migrate: ## Run database migrations
	go run functions/migration/main.go

migrate.atlas: ## Testing migration with atlas
	atlas schema apply --env gorm -u "mysql://root:password@localhost:3306/maindb"

migrate.undo: ## Undo the last database migration
	go run functions/migration/main.go --down

seed: ## Run database migrations
	go run functions/seed/main.go

test: ## Run tests
	scripts/test.sh

test.cover: test ## Run tests and open coverage statistics page
	go tool cover -html=coverage-all.out

clean: ## Clean up the built & test files
	rm -rf deploy/.serverless

specs: ## Generate swagger specs
	SWAGGER=$(SWAGGER) scripts/specs-gen.sh

specs.s: ## Generate swagger specs for swaggo
	SWAGGER=$(SWAGGER) scripts/specs-gen-swaggo.sh

build.api: ## Build the api services
	scripts/build-api.sh

build.funcs: ## Build the functions
	scripts/build-funcs.sh

build.api.lambda: ## Build the api services for AWS Lambda
	@$(MAKE) TARGET=lambda build.api

build.funcs.lambda: ## Build the functions for AWS Lambda
	@$(MAKE) TARGET=lambda build.funcs

deploy.dev: ## Deploy to DEV environment
	aws-vault exec [REPLACE_BY_AWS_PROFILE] -- scripts/deploy.sh dev

deploy.uat: ## Deploy to UAT environment
	aws-vault exec [REPLACE_BY_AWS_PROFILE] --no-session -- scripts/deploy.sh uat

deploy.prod: ## Deploy to PROD environment
	aws-vault exec [REPLACE_BY_AWS_PROFILE] --no-session -- scripts/deploy.sh prod
