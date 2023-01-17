LOCAL_DB_NAME:=slim_fairy_local
LOCAL_DB_DSN:=host=127.0.0.1 port=54320 dbname=slim_fairy_local user=slim_fairy_user password=slim_fairy_password sslmode=disable

docker-up-local:
	docker-compose --env-file ./configs/docker_local.env up -d --build
	goose -dir db/migrations postgres "$(LOCAL_DB_DSN)" up

docker-down-local:
	docker-compose down

db-create-migration:
	goose -dir db/migrations create name sql
