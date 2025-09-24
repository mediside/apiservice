MIGRATION_DIR = "./migrations"
DB_DRIVER = "postgres"
DB_STRING = "postgres://postgres:postgres@localhost:9432/med?sslmode=disable"

proto_dir         = internal/proto
proto_build_dir   = internal/gen/go/inference


# 
# 
compile-proto:
	protoc -I$(proto_dir) \
	--proto_path=$(proto_dir) \
	--go_out=$(proto_build_dir) \
	--go-grpc_out=$(proto_build_dir) \
	$(proto_dir)/*.proto

migrate-create: # make migrate-create name=init
	goose -dir ${MIGRATION_DIR} create $(name) sql

migrate-up:
	goose -dir ${MIGRATION_DIR} ${DB_DRIVER} ${DB_STRING} up

migrate-down:
	goose -dir ${MIGRATION_DIR} ${DB_DRIVER} ${DB_STRING} down

run:
	go run cmd/apiservice/main.go
