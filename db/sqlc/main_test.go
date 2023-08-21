package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbdriver = "postgres"
	dbsource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQuries Store
var testdb *sql.DB

func TestMain(m *testing.M) {
	var err error
	testdb, err = sql.Open(dbdriver, dbsource)
	if err != nil {
		log.Println(err)
	}
	testQuries = NewStore(testdb)
	os.Exit(m.Run())
}
