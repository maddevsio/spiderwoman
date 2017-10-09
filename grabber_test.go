package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

type ExternalServiceItem struct {
	Name string
	URL  string
}

type ExternalSeriviceList struct {
	Services []ExternalServiceItem
}

type Grabber interface {
	CheckConnection() (int, error)
	GetRawData() (string, error)
	// ParseData() (string, error)
}

func (d ExternalServiceItem) CheckConnection() (int, error) {
	resp, err := http.Get(d.URL)
	if err != nil {
		fmt.Println(err)
		return resp.StatusCode, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func (d ExternalServiceItem) GetRawData() (string, error) {
	resp, err := http.Get(d.URL)
	if err != nil {
		fmt.Println(err)
		return resp.Status, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(body), nil
}

func DefineExternalServices() ExternalSeriviceList {
	return ExternalSeriviceList{
		[]ExternalServiceItem{
			{"Alexa", "https://www.alexa.com"},
			{"Ahrefs", "https://www.ahrefs.com"},
		},
	}
}

func Grab(g Grabber) (int, error) {
	conn, err := g.CheckConnection()
	if err != nil {
		fmt.Println(err)
	}
	if conn == 200 {
		fmt.Println("Connection is ok, can proceed.")
		rawData, err := g.GetRawData()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(rawData)
		// TODO: Parse function is needs to be implemented for each External Service differently
	}
	return conn, err
}

func TestGrabber(t *testing.T) {
	defer gock.Off()
	gock.New("https://www.alexa.com").Get("/").Reply(200).BodyString("<h1>API</h1>")
	gock.New("https://www.alexa.com").Get("/").Reply(200).BodyString("<h1>API</h1>")
	gock.New("https://www.ahrefs.com").Get("/").Reply(200).BodyString("<h1>API</h1>")
	gock.New("https://www.ahrefs.com").Get("/").Reply(200).BodyString("<h1>API</h1>")

	grabberAPIURLs := DefineExternalServices()

	for _, service := range grabberAPIURLs.Services {
		fmt.Printf("Grabbing %v\n", service.URL)
		str, err := Grab(service)
		if err != nil {
			fmt.Println("fail", err)
		}
		assert.NoError(t, err)
		assert.Equal(t, 200, str)
	}
}
