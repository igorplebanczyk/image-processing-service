.PHONY: build run stop test lint sec

build:
	docker-compose build

run:
	docker-compose up

stop:
	docker-compose down

test:
	cd src && go test -v -cover ./...

lint:
	cd src && staticcheck ./...

sec:
	cd src && gosec ./...
