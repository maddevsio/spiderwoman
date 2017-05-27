package main

import (
	"testing"
	"gopkg.in/h2non/gock.v1"
	"os"
	"github.com/maddevsio/spiderwoman/lib"
	"github.com/stretchr/testify/assert"
)

func TestCrawl(t *testing.T) {
	defer gock.Off()

	dbPath := "./testdata/spiderwoman.db"

	os.Remove(dbPath)

	gock.New("http://server.com").
		Get("/").
		Reply(200).
		BodyString(
		"<a href='http://lalka.com'>mamka</a>" +
		"<a href='http://'>vse na mid!!</a>")

	path := Path{dbPath, "./testdata/sites.txt", "./sources.default.txt", "", "./types.default.txt"}
	initialize(path)
	crawl(path)

	monitors, err := lib.GetAllDataFromMonitor(dbPath, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(monitors))
	assert.Equal(t, "lalka.com", monitors[0].ExternalHost)
}

func TestCrawlCaseInsensitive(t *testing.T) {
	defer gock.Off()

	dbPath := "./testdata/spiderwoman.db"

	os.Remove(dbPath)

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

	path := Path{dbPath, "./testdata/sites.txt", "./sources.default.txt", "", "./types.default.txt"}
	initialize(path)
	crawl(path)

	monitors, err := lib.GetAllDataFromMonitor(dbPath, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(monitors))
	assert.Equal(t, "lalka.com", monitors[0].ExternalHost)
	assert.Equal(t, 2, monitors[0].Count)
}