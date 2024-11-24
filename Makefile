.PHONY: run stop test lint sec

run:
	docker-compose up --build

stop:
	docker-compose down

test:
	cd src && go test -v -cover ./...

lint:
	cd src && staticcheck ./...

sec:
	cd src && gosec ./...
