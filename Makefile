MIGRATION_DIR = "./migrations"
DB_DRIVER = "postgres"
DB_STRING = "postgres://postgres:postgres@localhost:9432/med?sslmode=disable"


migrate-create: # make migrate-create name=init
	goose -dir ${MIGRATION_DIR} create $(name) sql

migrate-up:
	goose -dir ${MIGRATION_DIR} ${DB_DRIVER} ${DB_STRING} up

migrate-down:
	goose -dir ${MIGRATION_DIR} ${DB_DRIVER} ${DB_STRING} down

run:
	go run cmd/apiservice/main.go
