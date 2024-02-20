.DEFAULT_GOAL := help

# HOST is only used for API specs generation
HOST ?= localhost:8080

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

start: docker.up ## Bring up the server on dev environment
	sh scripts/waitdb.sh
	air

docker.down: ## Bring down the server on dev environment, remove all docker related stuffs as well
	docker-compose down -v --remove-orphans

docker.up: ## Bring up the docker container
	docker-compose up -d

remove: ## Bring down the server on dev environment, remove all docker related stuffs as well
	docker-compose down -v --remove-orphans

migrate: ## Run database migrations
	go run cmd/migration/main.go

migrate.undo: ## Undo the last database migration
	go run cmd/migration/main.go --down

seed: ## Run database seeder
	echo "To be done!"

test: ## Run tests
	sh scripts/test.sh

test.cover: test ## Run tests and open coverage statistics page
	go tool cover -html=coverage-all.out

build: clean ## Build the server binary file on host machine
	sh scripts/build.sh

build.linux: ## Build the server binary file for Linux host
	@$(MAKE) GOOS=linux GOARCH=amd64 build

build.windows: ## Build the server binary file for Windows host
	@$(MAKE) GOOS=windows GOARCH=amd64 build

build.arm: clean ## Build the server binary file for ARM host
	GOOS=linux GOARCH=arm64 sh scripts/build-arm.sh

install:
	echo "Not ready yet!"
	echo "To setup PostgreSQL, check 'scripts/install-pg.sh'"
	echo "To setup the server, check 'scripts/install-service.sh'"

clean: ## Clean up the built & test files
	rm -rf ./server ./*.out

specs: ## Generate swagger specs
	HOST=$(HOST) sh scripts/specs-gen.sh

%: # prevent error for `up` target when passing arguments
ifeq ($(filter up,$(MAKECMDGOALS)),up)
	@:
else
	$(error No rule to make target `$@`.)
endif
