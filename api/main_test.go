package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"os"
	"strconv"
	"github.com/maddevsio/spiderwoman/lib"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/maddevsio/simple-config"
)

var (
	config = simple_config.NewSimpleConfig("../config.test", "yml")
)

func TestIndex(t *testing.T) {
	ts := httptest.NewServer(GetAPIEngine(config))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, string(actual), "Spiderwoman")
}

func TestPing200(t *testing.T) {
	ts := httptest.NewServer(GetAPIEngine(config))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}

func TestPing404(t *testing.T) {
	ts := httptest.NewServer(GetAPIEngine(config))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/pong")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 404 {
		t.Fatalf("Received non-404 response: %d\n", resp.StatusCode)
	}
}

func TestAll(t *testing.T) {
	os.Remove(config.GetString("db-path"))
	lib.CreateDBIfNotExists(config.GetString("db-path"))

	for i := int(0); i < 2; i++ {
		sourceHost := "http://a"
		externalLink := "http://b/1?" + strconv.Itoa(i)
		count := 800+i
		externalHost := "b"
		_ = lib.SaveRecordToMonitor(config.GetString("db-path"), sourceHost, externalLink, count, externalHost)
	}

	ts := httptest.NewServer(GetAPIEngine(config))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/all")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var r []lib.Monitor
	err = json.Unmarshal([]byte(actual), &r)
	if err != nil {
		t.Fatal(err)
	}
}