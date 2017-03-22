package main

import (
	"testing"
	"gopkg.in/h2non/gock.v1"
	"os"
	"github.com/maddevsio/spiderwoman/lib"
	"github.com/stretchr/testify/assert"
)

func initMockHTTPServer() {
	// TODO: need to try
	// https://github.com/h2non/gock
	// http://www.mock-server.com/
	// https://github.com/jarcoal/httpmock
	// https://github.com/golang/mock
}

func TestCrawl(t *testing.T) {
	defer gock.Off()

	dbPath := "./testdata/spiderwoman.db"

	os.Remove(dbPath)

	gock.New("http://server.com").
		Get("/").
		Reply(200).
		BodyString(
		"<a href='http://lalka.com'>blah</a>" +
		"<a href='http://'>blah</a>")

	path := Path{dbPath, "./testdata/sites.txt", "./sites.default.txt"}
	initialize(path)
	crawl(path)

	monitors, err := lib.GetAllDataFromMonitor(dbPath, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(monitors))
}