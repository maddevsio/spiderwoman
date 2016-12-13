package lib

import (
	"database/sql"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	DBFilepath = "/tmp/spiderwoman.db"
)

func TestCreateDBIfNotExists(t *testing.T) {
	os.Remove(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	_, err := os.Stat(DBFilepath)
	assert.Equal(t, false, os.IsNotExist(err))
}

func TestCheckMonitorTable(t *testing.T) {
	os.Remove(DBFilepath)
	CreateDBIfNotExists(DBFilepath)

	db, err := sql.Open("sqlite3", DBFilepath)
	assert.Equal(t, nil, err)
	defer db.Close()

	rows, err := db.Query("SELECT name FROM sqlite_master;")
	assert.Equal(t, nil, err)
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		assert.Equal(t, nil, err)
		assert.Equal(t, "monitor", name)
	}

	err = rows.Err()
	assert.Equal(t, nil, err)
}

func TestSaveRecordToMonitor(t *testing.T) {
	source_host := "http://a"
	external_link := "http://b/1"
	count := 800
	external_host := "b"
	res := SaveRecordToMonitor(DBFilepath, source_host, external_link, count, external_host)
	assert.Equal(t, true, res)

	db, err := sql.Open("sqlite3", DBFilepath)
	assert.Equal(t, nil, err)
	defer db.Close()

	rows, err := db.Query("SELECT source_host, created FROM monitor;")
	assert.Equal(t, nil, err)
	defer rows.Close()

	k := 0
	for rows.Next() {
		k += 1
		var source_host string
		var created string
		err = rows.Scan(&source_host, &created)
		assert.Equal(t, nil, err)
		assert.Equal(t, "http://a", source_host)
	}
	assert.Equal(t, 1, k)
}

func TestParseSqliteDate(t *testing.T) {
	sqliteDate := "2016-12-13T07:17:23Z"
	parsedTime, _ := ParseSqliteDate(sqliteDate)
	assert.Equal(t, "2016", strconv.Itoa(parsedTime.Year()))
	assert.Equal(t, "12", strconv.Itoa(int(parsedTime.Month())))
	assert.Equal(t, "13", strconv.Itoa(parsedTime.Day()))
}
