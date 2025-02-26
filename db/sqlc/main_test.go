package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"os"
	"simple-banking/util"
	"testing"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatalf("cannot load config file: %v", err)
	}

	testDb, err = sql.Open(config.DbDriver, config.DataSource)
	if err != nil {
		log.Fatal("can not connect to database")
	}
	testQueries = New(testDb)
	os.Exit(m.Run())
}
