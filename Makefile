.PHONY: build up down restart logs tidy test lint sec tf
.DEFAULT_GOAL := up

SRC_DIR := src
TF_DIR := terraform

# Docker Compose commands

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

logs:
	docker-compose logs -f

# Go commands

tidy:
	cd $(SRC_DIR); go mod tidy

test:
	cd $(SRC_DIR); go test -v -cover ./...

lint:
	cd $(SRC_DIR); staticcheck ./...

sec:
	cd $(SRC_DIR); gosec ./...

# Terraform commands; usage: make tf <command> (e.g. make tf apply)

tf:
	cd $(TF_DIR); terraform $(filter-out $@,$(MAKECMDGOALS))
