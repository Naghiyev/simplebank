package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"simple-banking/api"
	db "simple-banking/db/sqlc"
	"simple-banking/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DbDriver, config.DataSource)

	if err != nil {
		log.Fatal("can not connect to db", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Server couldn't started", err)
	}

}
