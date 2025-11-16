COMPOSE=docker-compose
.PHONY: all up build down clean rebuild

all: up

up:
	$(COMPOSE) up --build

build:
	$(COMPOSE) build

down:
	$(COMPOSE) down

test:
	go test ./...

clean:
	$(COMPOSE) down -v

rebuild: down build up