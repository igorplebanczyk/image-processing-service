.PHONY: build run

OUTPUT=ips

build:
	cd cmd && go build -o ../$(OUTPUT)

run:
	./$(OUTPUT)

all: build run