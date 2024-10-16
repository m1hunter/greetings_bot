ENV_FILE=.env

TELEGRAM_BOT_TOKEN ?= TELEGRAM_BOT_TOKEN=7090855604:AAE9Ui7J34yFjnQ7Be8EyRtomfkxTVFm7CU
DB_HOST ?= db
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= postgres

$(ENV_FILE):
	@echo "TELEGRAM_BOT_TOKEN=$(TELEGRAM_BOT_TOKEN)" > $(ENV_FILE)
	@echo "DB_HOST=$(DB_HOST)" >> $(ENV_FILE)
	@echo "DB_PORT=$(DB_PORT)" >> $(ENV_FILE)
	@echo "DB_USER=$(DB_USER)" >> $(ENV_FILE)
	@echo "DB_PASSWORD=$(DB_PASSWORD)" >> $(ENV_FILE)
	@echo "DB_NAME=$(DB_NAME)" >> $(ENV_FILE)

up: $(ENV_FILE)
	@docker-compose up --build

down:
	@docker-compose down

clean: down
	@rm -f $(ENV_FILE)

default: up
