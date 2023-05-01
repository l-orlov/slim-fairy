LOCAL_DB_HOST:=127.0.0.1
LOCAL_DB_NAME:=slim_fairy_local
LOCAL_DB_USER:=slim_fairy_user
LOCAL_DB_PASSWORD:=slim_fairy_password
LOCAL_DB_PORT:=54320
LOCAL_DB_DSN:=host=$(LOCAL_DB_HOST) port=$(LOCAL_DB_PORT) dbname=$(LOCAL_DB_NAME) user=$(LOCAL_DB_USER) password=$(LOCAL_DB_PASSWORD) sslmode=disable

# set up locally with docker
docker-up-local:
	docker-compose --env-file ./configs/docker_local.env up -d --build
	# sleep before run migrations to wait db creation
	sleep 1
	goose -dir db/migrations postgres "$(LOCAL_DB_DSN)" up

docker-down-local:
	docker-compose --env-file ./configs/docker_local.env down

docker-reset-local:
	make docker-down-local
	make docker-up-local

db-create-migration:
	goose -dir db/migrations create name sql

db-migrate:
	goose -dir db/migrations postgres "$(LOCAL_DB_DSN)" up
	make db-gen-structure

db-migrate-down:
	goose -dir db/migrations postgres "$(LOCAL_DB_DSN)" down
	make db-gen-structure

db-gen-structure:
	pg_dump "$(LOCAL_DB_DSN)" --schema-only --no-owner --no-privileges --no-tablespaces --no-security-labels --no-comments > db/structure.sql

# for local db testing
db-reset:
	psql -c "drop database if exists $(LOCAL_DB_NAME)"
	psql -c "create database $(LOCAL_DB_NAME)"
	goose -dir db/migrations postgres "$(LOCAL_DB_DSN)" up
	make db-gen-structure

# generate project code
gen:
	go generate ./internal/...

run:
	go run cmd/slim-fairy/main.go
