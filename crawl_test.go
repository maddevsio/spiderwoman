package main

import "testing"

func initMockHTTPServer() {}

func TestCrawl(t *testing.T) {
	initialize()
	crawl()
	// TODO: pass struct to initialize and crawl
	// TODO: remove config dependencies
}