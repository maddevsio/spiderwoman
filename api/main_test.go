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
)

const (
	DBFilepath = "/tmp/spiderwoman.db"
)

func TestPing200(t *testing.T) {
	ts := httptest.NewServer(GetAPIEngine(DBFilepath))
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
	ts := httptest.NewServer(GetAPIEngine(DBFilepath))
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
	os.Remove(DBFilepath)
	lib.CreateDBIfNotExists(DBFilepath)

	for i := int(0); i < 2; i++ {
		sourceHost := "http://a"
		externalLink := "http://b/1?" + strconv.Itoa(i)
		count := 800+i
		externalHost := "b"
		_ = lib.SaveRecordToMonitor(DBFilepath, sourceHost, externalLink, count, externalHost)
	}

	ts := httptest.NewServer(GetAPIEngine(DBFilepath))
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