DB_URL := "postgres://sdai-calculator:sdai-calculator@localhost:5432/sdai-calculator?sslmode=disable"
MIGRATIONS_PATH := storage/postgresql/migrations

docker-up:
	docker compose up -d

docker-down:
	docker compose down

run:
	@make docker-up
	go run cmd/sdai-calculator/main.go --config=./config/local.yaml

mig:
	go run ./cmd/migrator/ --db-url=${DB_URL} --migrations-path="file://${MIGRATIONS_PATH}"

gen:
	mkdir -p ./internal/generated ./internal/generated/server
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=./api/server/config.yaml ./api/server/openapi.yaml

clean:
	rm -rf ./internal/generated

lint:
	golangci-lint run ./...

drop:
	migrate -database ${DB_URL} -path "${MIGRATIONS_PATH}" drop
