package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
)

func TestPing200(t *testing.T) {
	ts := httptest.NewServer(GetAPIEngine())
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
	ts := httptest.NewServer(GetAPIEngine())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/pong")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 404 {
		t.Fatalf("Received non-404 response: %d\n", resp.StatusCode)
	}
}