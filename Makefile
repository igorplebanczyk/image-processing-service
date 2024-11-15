.PHONY: build run stop test coverage

build:
	docker-compose build

run:
	docker-compose up

stop:
	docker-compose down

test:
	cd src && go test -v ./...

coverage:
	cd src && go test -cover ./...