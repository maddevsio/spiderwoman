package main

import (
	"testing"
	"gopkg.in/h2non/gock.v1"
	"github.com/maddevsio/spiderwoman/lib"
	"github.com/stretchr/testify/assert"
	"log"
)

func TestCrawl(t *testing.T) {
	defer gock.Off()

	dbName := "spiderwoman-test-crawl"

	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExists(dbName)

	gock.New("http://server.com").
		Get("/").
		Reply(200).
		BodyString(
		"<a href='http://lalka.com'>mamka</a>" +
		"<a href='http://'>vse na mid!!</a>")

	path := Path{dbName, "./testdata/sites.txt", "./sources.default.txt", "", "./types.default.txt"}
	initialize(path)
	crawl(path)

	monitors, err := lib.GetAllDataFromMonitor(dbName, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(monitors))
	assert.Equal(t, "lalka.com", monitors[0].ExternalHost)
}

func TestCrawlCaseInsensitive(t *testing.T) {
	defer gock.Off()

	dbName := "spiderwoman-test-crawl"

	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExists(dbName)

	gock.New("http://server.com").
		Get("/").
		Reply(200).
		BodyString(
		"<a href='http://lalka.com'>mamka</a>" +
		"<a href='/2'>inteernal</a>" +
		"<a href='http://'>vse na mid!!</a>")

	gock.New("http://server.com").
		Get("/2").
		Reply(200).
		BodyString("<a href='http://laLka.com'>maMka</a>")

	path := Path{dbName, "./testdata/sites.txt", "./sources.default.txt", "", "./types.default.txt"}
	initialize(path)
	crawl(path)

	monitors, err := lib.GetAllDataFromMonitor(dbName, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(monitors))
	assert.Equal(t, "lalka.com", monitors[0].ExternalHost)
	assert.Equal(t, 2, monitors[0].Count)
}

func TestCrawlCase1(t *testing.T) {
	defer gock.Off()

	dbName := "spiderwoman-test-crawl"

	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExists(dbName)

	gock.New("http://server.com").
		Get("/").
		Reply(200).
		BodyString(
		"<a href='http://lalka.com'>mamka</a>" +
			"<a href='http://LALkA.cOm'>azzazaza</a>")

	path := Path{dbName, "./testdata/sites.txt", "./sources.default.txt", "", "./types.default.txt"}
	initialize(path)
	crawl(path)

	monitors, err := lib.GetAllDataFromMonitor(dbName, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(monitors))
	assert.Equal(t, "lalka.com", monitors[0].ExternalHost)
}

func TestCrawlCase2(t *testing.T) {
	// TODO: gock is not properly working when http call goes from goroutines. Need to find new mock server or ask the community

	defer gock.Off()

	dbName := "spiderwoman-test-crawl"

	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExists(dbName)

	gock.New("http://server.com").
		Get("/").
		Persist().
		Reply(200).
		BodyString("<a href='http://server.com/go/1'>redirect</a> | <a href='http://lalka.com/lool'>regular link</a>")

	gock.New("http://server.com").
		Get("/go/1").
		Persist().
		Reply(302).
		SetHeader("Location", "https://www.google.com")

	gock.New("http://lalka.com").
		Get("/lool").
		Reply(200).
		BodyString("plaaaah")

	path := Path{dbName, "./testdata/sites.txt", "./sources.default.txt", "", "./types.default.txt"}
	initialize(path)
	crawl(path)

	monitors, err := lib.GetAllDataFromMonitor(dbName, 0)

	log.Printf("%v", monitors)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(monitors))
	// assert.Equal(t, "lalka.com", monitors[0].ExternalHost)
	// TODO: FIX THIS

}

func TestCrawlWWW(t *testing.T) {
	defer gock.Off()

	dbName := "spiderwoman-test-crawl"

	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExists(dbName)

	gock.New("http://server.com").
		Get("/").
		Reply(200).
		BodyString(
		"<a href='http://lalka.com'>mamka</a>" +
			"<a href='http://www.lalka.COM/blaah'>mamka</a>")

	path := Path{dbName, "./testdata/sites.txt", "./sources.default.txt", "", "./types.default.txt"}
	initialize(path)
	crawl(path)

	monitors, err := lib.GetAllDataFromMonitor(dbName, 0)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(monitors))
	assert.Equal(t, "lalka.com", monitors[0].ExternalHost)
	assert.Equal(t, "lalka.com", monitors[1].ExternalHost)
}

func TestCrawl443(t *testing.T) {
	defer gock.Off()

	dbName := "spiderwoman-test-crawl"

	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExists(dbName)

	gock.New("http://server.com").
		Get("/").
		Reply(200).
		BodyString(
		"<a href='http://lalka.com'>mamka</a>" +
			"<a href='http://www.lalka.COM:443/blaah'>mamka</a>")

	path := Path{dbName, "./testdata/sites.txt", "./sources.default.txt", "", "./types.default.txt"}
	initialize(path)
	crawl(path)

	monitors, err := lib.GetAllDataFromMonitor(dbName, 0)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(monitors))
	assert.Equal(t, "lalka.com", monitors[0].ExternalHost)
	assert.Equal(t, "lalka.com", monitors[1].ExternalHost)
}


func TestCrawUnicodeLink(t *testing.T) {
	defer gock.Off()

	dbName := "spiderwoman-test-crawl"

	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExists(dbName)

	gock.New("http://server.com").
		Get("/").
		Reply(200).
		BodyString("<a href='https://���������.com/?i=382'>mamka</a>")

	path := Path{dbName, "./testdata/sites.txt", "./sources.default.txt", "", "./types.default.txt"}
	initialize(path)
	crawl(path)

	monitors, err := lib.GetAllDataFromMonitor(dbName, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(monitors))
	assert.Equal(t, "���������.com", monitors[0].ExternalHost)
}