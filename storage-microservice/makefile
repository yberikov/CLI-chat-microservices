
migrate-up:
	@echo "Running migrate command..."
	migrate -path ./schema -database 'postgres://user:12345@localhost:5434/chat?sslmode=disable' up
migrate-down:
	migrate -path ./schema -database 'postgres://postgres:12345@localhost:5434/chat?sslmode=disable' down

.PHONY: migrate-up migrate-down