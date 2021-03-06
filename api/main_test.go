package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/maddevsio/simple-config"
	"github.com/maddevsio/spiderwoman/lib"
	"github.com/stretchr/testify/assert"
	"github.com/tealeg/xlsx"
)

var (
	config = simple_config.NewSimpleConfig("../config.test", "yml")
)

func TestIndex(t *testing.T) {
	lib.TruncateDB(config.GetString("db-path"))
	lib.CreateDBIfNotExistsAndMigrate(config.GetString("db-path"))

	ts := httptest.NewServer(GetAPIEngine(config))
	defer ts.Close()

	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL+"/", nil)
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
	req, err := http.NewRequest("GET", ts.URL+"/ping", nil)
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
	req, err := http.NewRequest("GET", ts.URL+"/pong", nil)
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
	lib.TruncateDB(config.GetString("db-path"))
	lib.CreateDBIfNotExistsAndMigrate(config.GetString("db-path"))

	for i := int(0); i < 2; i++ {
		m := lib.Monitor{}
		m.SourceHost = "http://a"
		m.ExternalLink = "http://b/1?" + strconv.Itoa(i)
		m.Count = 800 + i
		m.ExternalHost = "b"
		_ = lib.SaveRecordToMonitor(config.GetString("db-path"), m)
	}

	ts := httptest.NewServer(GetAPIEngine(config))
	defer ts.Close()

	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL+"/all", nil)
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

func TestXLSSimple(t *testing.T) {
	xlsTestFile := "/tmp/test-excel.xls"
	uris := [2]string{"/get-new-xls", "/get-day-xls"}
	for _, uri := range uris {
		os.Remove(xlsTestFile)
		ts := httptest.NewServer(GetAPIEngine(config))
		defer ts.Close()

		client := &http.Client{}
		req, err := http.NewRequest("GET", ts.URL+uri, nil)
		req.SetBasicAuth(config.GetString("admin-user"), config.GetString("admin-password"))
		resp, err := client.Do(req)

		if err != nil {
			t.Fatal(err)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		err = ioutil.WriteFile(xlsTestFile, data, 0644)
		assert.NoError(t, err)

		excelFileName := xlsTestFile
		_, err = xlsx.OpenFile(excelFileName)
		assert.NoError(t, err)
	}
}

func TestAllForHost(t *testing.T) {
	lib.TruncateDB(config.GetString("db-path"))
	lib.CreateDBIfNotExistsAndMigrate(config.GetString("db-path"))

	for i := int(0); i < 2; i++ {
		m := lib.Monitor{}
		m.SourceHost = "http://a"
		m.ExternalLink = "http://b/1?" + strconv.Itoa(i)
		m.Count = 800 + i
		m.ExternalHost = "b"
		_ = lib.SaveRecordToMonitor(config.GetString("db-path"), m)
	}

	ts := httptest.NewServer(GetAPIEngine(config))
	defer ts.Close()

	// check without "host" query string param

	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL+"/all-for-host", nil)
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
	req, err = http.NewRequest("GET", ts.URL+"/all-for-host?host=b", nil)
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
