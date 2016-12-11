package lib

import (
	"database/sql"
	"log"
	"os"
	"testing"
)

const (
	DBFilepath = "/tmp/spiderwoman.db"
)

func TestCreateDBIfNotExists(t *testing.T) {
	os.Remove(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	if _, err := os.Stat(DBFilepath); os.IsNotExist(err) {
		t.Error("DB file does not exist")
	}
}

func TestCheckMonitorTable(t *testing.T) {
	os.Remove(DBFilepath)
	CreateDBIfNotExists(DBFilepath)

	db, err := sql.Open("sqlite3", DBFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT name FROM sqlite_master;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		if name != "monitor" {
			t.Error("table MONITOR not found")
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func TestSaveRecordToMonitor(t *testing.T) {
	source_host := "http://a"
	external_link := "http://b/1"
	count := 800
	external_host := "b"
	SaveRecordToMonitor(DBFilepath, source_host, external_link, count, external_host)
}
