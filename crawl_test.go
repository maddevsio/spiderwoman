package main

import "testing"

func initMockHTTPServer() {
	// TODO: need to try
	// https://github.com/h2non/gock
	// http://www.mock-server.com/
	// https://github.com/jarcoal/httpmock
	// https://github.com/golang/mock
}

func TestCrawl(t *testing.T) {
	path := Path{sqliteDBPath, "./testdata/sites.txt", "./sites.default.txt"}
	initialize(path)
	crawl(path)
}