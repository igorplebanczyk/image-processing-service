.PHONY: run stop test lint sec tfinit tfplan tfapply

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

tfinit:
	cd terraform && terraform init

tfplan:
	cd terraform && terraform plan

tfapply:
	cd terraform && terraform apply
