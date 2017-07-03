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

	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL + "/", nil)
	req.SetBasicAuth(config.GetString("admin-user"), config.GetString("admin-password"))
	resp, err := client.Do(req)
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

	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL + "/ping", nil)
	req.SetBasicAuth(config.GetString("admin-user"), config.GetString("admin-password"))
	resp, err := client.Do(req)

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

	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL + "/pong", nil)
	req.SetBasicAuth(config.GetString("admin-user"), config.GetString("admin-password"))
	resp, err := client.Do(req)

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

	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL + "/all", nil)
	req.SetBasicAuth(config.GetString("admin-user"), config.GetString("admin-password"))
	resp, err := client.Do(req)

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

func TestAllForDate(t *testing.T) {

}

func TestXLS(t *testing.T) {
	// download xls (can be empty)
	// read xls via https://github.com/tealeg/xlsx
	// check for errors
}

func TestAllForHost(t *testing.T) {
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

	// check without "host" query string param

	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL + "/all-for-host", nil)
	req.SetBasicAuth(config.GetString("admin-user"), config.GetString("admin-password"))
	resp, err := client.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var r1 []lib.Monitor
	err = json.Unmarshal([]byte(actual), &r1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, r1, 0)

	// check with "host" query string

	client = &http.Client{}
	req, err = http.NewRequest("GET", ts.URL + "/all-for-host?host=b", nil)
	req.SetBasicAuth(config.GetString("admin-user"), config.GetString("admin-password"))
	resp, err = client.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	actual, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var r2 []lib.Monitor
	err = json.Unmarshal([]byte(actual), &r2)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, r2, 2)
}