PROJECT = biocad
COMPOSE_FILES = -f docker-compose.yml
COMPOSE = docker-compose -p $(PROJECT) $(COMPOSE_FILES)

include .env
export $(shell sed 's/=.*//' .env)

all:
	$(COMPOSE) up --build

build:
	$(COMPOSE) build

db:
	$(COMPOSE) up -d biocad_psql

rebuild:
	$(COMPOSE) up --build --force-recreate -d server

clean:
	$(COMPOSE) down -v --rmi local