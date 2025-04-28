MIGRATE_CMD = go run cmd/migrate/main.go

migrate:
	@$(MIGRATE_CMD)

migrate-down:
	@$(MIGRATE_CMD) down

migrate-status:
	@$(MIGRATE_CMD) status

sqlc-generate:
	sqlc generate

swagger:
	swag init -o internal/server/docs -g cmd/server/main.go

#BUILDING
APP_NAME=go-web
DOCKER_CONTAINER=$(APP_NAME)-container
PORT=3000

build:
	docker build -t $(APP_NAME) .

run:
	docker run --rm --name $(DOCKER_CONTAINER) --env-file .env -p $(PORT):$(PORT) $(APP_NAME)

clean:
	-docker stop $(DOCKER_CONTAINER)
	-docker rm $(DOCKER_CONTAINER)
