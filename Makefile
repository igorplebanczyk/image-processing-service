.PHONY: build run stop test

build:
	docker-compose build

run:
	docker-compose up

stop:
	docker-compose down

test:
	go test -v ./...