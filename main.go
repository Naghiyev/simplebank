package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"simple-banking/api"
	db "simple-banking/db/sqlc"
)

const (
	dbdriver       = "postgres"
	dataSourceName = "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable"
	Serveraddress  = "localhost:7070"
)

func main() {
	conn, err := sql.Open(dbdriver, dataSourceName)

	if err != nil {
		log.Fatal("can not connect to db", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(Serveraddress)
	if err != nil {
		log.Fatal("Server couldn't started", err)
	}

}
