.PHONY: build-run run stop

build-run:
	docker-compose up --build

run:
	docker-compose up

stop:
	docker-compose down