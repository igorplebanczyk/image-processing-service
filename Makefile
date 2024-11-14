.PHONY: build run stop test coverage

build:
	docker-compose build

run:
	docker-compose up

stop:
	docker-compose down

test:
	go test -v ./...

coverage:
	go test -cover ./...