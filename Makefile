db:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16.2-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root go_chat

drop_db:
	docker exec -it postgres dropdb go_chat

db_cli:
	docker exec -it postgres psql -U root

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/go_chat?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/go_chat?sslmode=disable" -verbose down

sqlc:
	sqlc generate

server:
	go run main.go

create_migrations:
	migrate create -ext sql -dir db/migrations add_users_table

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine


.PHONY: db createdb drop_db db_cli redis