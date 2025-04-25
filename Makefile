MIGRATE_CMD = go run cmd/migrate/main.go

#Database migrations
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
