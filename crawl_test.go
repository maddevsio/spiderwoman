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
	initialize()
	crawl()
	// TODO: pass struct to initialize and crawl
	// TODO: remove config dependencies

}