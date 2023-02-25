LOCAL_DB_NAME:=slim_fairy_local
LOCAL_DB_DSN:=host=127.0.0.1 port=54320 dbname=$(LOCAL_DB_NAME) user=slim_fairy_user password=slim_fairy_password sslmode=disable

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


# generate project code
gen:
	go generate ./internal/...

run:
	go run cmd/slim-fairy/main.go
