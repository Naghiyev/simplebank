postgres:
	docker run --name postgres-simple-banking  -e POSTGRES_USER=root  -e POSTGRES_PASSWORD=root   -p 5432:5432 -d postgres

createdb:
	docker exec -it postgres-simple-banking createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres-simple-banking dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go simple-banking/db/sqlc Store


.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock
