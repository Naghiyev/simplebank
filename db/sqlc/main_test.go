package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

const (
	dbdriver       = "postgres"
	dataSourceName = "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDb, err = sql.Open(dbdriver, dataSourceName)
	if err != nil {
		log.Fatal("can not connect to database")
	}
	testQueries = New(testDb)
	os.Exit(m.Run())
}
