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
	os.Remove(DBFilepath)
	CreateDBIfNotExists(DBFilepath)

	sourceHost := "http://a"
	externalLink := "http://b/1"
	count := 800
	externalHost := "b"
	res := SaveRecordToMonitor(DBFilepath, sourceHost, externalLink, count, externalHost)
	assert.Equal(t, true, res)

	db, err := sql.Open("sqlite3", DBFilepath)
	assert.Equal(t, nil, err)
	defer db.Close()

	rows, err := db.Query("SELECT source_host, created FROM monitor;")
	assert.Equal(t, nil, err)
	defer rows.Close()

	k := 0
	for rows.Next() {
		k++
		var sourceHost string
		var created string
		err = rows.Scan(&sourceHost, &created)
		assert.Equal(t, nil, err)
		assert.Equal(t, "http://a", sourceHost)
	}
	assert.Equal(t, 1, k)
}

func TestSaveRecordToMonitor_BadExternalLink(t *testing.T) {
	os.Remove(DBFilepath)
	CreateDBIfNotExists(DBFilepath)

	sourceHost := "http://a"
	externalLink := "http://somebaddomain.com/'.show_site_name($line['url'],100).'/"
	count := 800
	externalHost := "b"
	res := SaveRecordToMonitor(DBFilepath, sourceHost, externalLink, count, externalHost)
	assert.Equal(t, true, res)

	db, err := sql.Open("sqlite3", DBFilepath)
	assert.Equal(t, nil, err)
	defer db.Close()

	rows, err := db.Query("SELECT external_link FROM monitor;")
	assert.Equal(t, nil, err)
	defer rows.Close()

	k := 0
	for rows.Next() {
		k++
		var fetchedExternalLink string
		err = rows.Scan(&fetchedExternalLink)
		assert.Equal(t, nil, err)
		assert.Equal(t, externalLink, fetchedExternalLink)
	}
	assert.Equal(t, 1, k)
}

func TestGetAllDataFromSqlite_MapToStruct(t *testing.T) {
	os.Remove(DBFilepath)
	CreateDBIfNotExists(DBFilepath)

	for i := int(0); i < 10; i++ {
		sourceHost := "http://a"
		externalLink := "http://b/1?" + strconv.Itoa(i)
		count := 800+i
		externalHost := "b"
		_ = SaveRecordToMonitor(DBFilepath, sourceHost, externalLink, count, externalHost)
	}

	monitors, _ := GetAllDataFromMonitor(DBFilepath)
	assert.Equal(t, 10, len(monitors))
	assert.Equal(t, "http://b/1?0", monitors[0].ExternalLink)
	assert.Equal(t, "http://b/1?9", monitors[9].ExternalLink)
}

func TestParseSqliteDate(t *testing.T) {
	sqliteDate := "2016-12-13T07:17:23Z"
	parsedTime, _ := ParseSqliteDate(sqliteDate)
	assert.Equal(t, "2016", strconv.Itoa(parsedTime.Year()))
	assert.Equal(t, "12", strconv.Itoa(int(parsedTime.Month())))
	assert.Equal(t, "13", strconv.Itoa(parsedTime.Day()))
}

func TestCrawlStatus(t *testing.T) {
	//SetCrawlStatus(DBFilepath, "Crawling...")
	//GetCrawlStatus(DBFilepath)
	//SetCrawlStatus(DBFilepath, "Crawl done")
	//GetCrawlStatus(DBFilepath)
}

