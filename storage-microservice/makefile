
run-postgres:
	docker run -d \
      --name postgres \
      -e POSTGRES_PASSWORD=12345 \
      -p 5436:5432 \
      -d --rm \
      postgres:latest


migrate-up:
	@echo "Running migrate command..."
	migrate -path ./schema -database 'postgres://postgres:12345@localhost:5433/postgres?sslmode=disable' up
migrate-down:
	migrate -path ./schema -database 'postgres://postgres:12345@localhost:5433/postgres?sslmode=disable' down

.PHONY: migrate-up migrate-down