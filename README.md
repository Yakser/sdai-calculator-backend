# sdai-calculator-backend


## Start service

- run codegen: `make gen`
- up infra in docker: `make docker-up`
- and then start server: `make run`.

## Migrations

Make sure that database container is running: `make docker-up`

Create migration:

```shell

migrate create -ext sql -dir storage/postgresql/migrations -seq create_users_table

```

Run migration:

```shell

migrate -database "postgres://sdai-calculator:sdai-calculator@localhost:5432/sdai-calculator?sslmode=disable" -path "storage/postgresql/migrations" up

```