DOCKER_DIR := docker
SRC_DIR := src
TF_DIR := terraform

.PHONY: help
help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

# Docker Compose commands

.ONESHELL:
.PHONY: build
build: ## Build the Docker containers
	cd $(DOCKER_DIR);
	docker-compose build

.ONESHELL:
.PHONY: up
up: ## Start the Docker containers
	cd $(DOCKER_DIR);
	docker-compose up -d

.ONESHELL:
.PHONY: down
down: ## Stop the Docker containers
	cd $(DOCKER_DIR);
	docker-compose down

.ONESHELL:
.PHONY: restart
restart: ## Restart the Docker containers
	cd $(DOCKER_DIR);
	docker-compose restart

.ONESHELL:
.PHONY: logs
logs: ## Show the logs of the Docker containers
	cd $(DOCKER_DIR);
	docker-compose logs -f

# Go commands

.ONESHELL:
.PHONY: init
init: ## Install development tools and project dependencies
	go install github.com/securego/gosec/v2/cmd/gosec@latest;
	go install honnef.co/go/tools/cmd/staticcheck@latest;
	cd $(SRC_DIR);
	go mod download

.ONESHELL:
.PHONY: tidy
tidy: ## Tidy Go dependencies
	cd $(SRC_DIR);
	go mod tidy

.ONESHELL:
.PHONY: test
test: ## Run Go tests
	cd $(SRC_DIR);
	go test -v -cover ./...

.ONESHELL:
.PHONY: format
format: ## Format Go code
	cd $(SRC_DIR);
	go fmt ./...

.ONESHELL:
.PHONY: vet
vet: ## Vet Go code
	cd $(SRC_DIR);
	go vet ./...

.ONESHELL:
.PHONY: lint
lint: ## Lint Go code
	cd $(SRC_DIR);
	staticcheck ./...

.ONESHELL:
.PHONY: sec
sec: ## Run Go security checks
	cd $(SRC_DIR);
	gosec ./...

# Terraform commands

.ONESHELL:
.PHONY: tf
tf: ## Run Terraform commands in the Terraform directory: make tf <command>; e.g. make tf plan
	cd $(TF_DIR);
	terraform $(filter-out $@,$(MAKECMDGOALS))
