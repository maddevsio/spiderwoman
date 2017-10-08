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
}

type InputData struct {
	FeaturedHost          string
	ExternalServiceAPIURL string
}

func (d InputData) CheckConnection() (int, error) {
	resp, err := http.Get(d.ExternalServiceAPIURL)
	if err != nil {
		fmt.Println(err)
		return resp.StatusCode, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func (d InputData) GetRawData() (string, error) {
	resp, err := http.Get(d.ExternalServiceAPIURL)
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
			{"Alexa", "http://alexa.com/api/"},
			{"Ahrefs", "https://ahrefs.com/api/v4/"},
			{"Google Metrics", "https://metrics.google.com/api/"},
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
	d := InputData{
		FeaturedHost:          "nambataxi.kg",
		ExternalServiceAPIURL: "https://www.alexa.com",
	}
	defer gock.Off()
	gock.New("https://www.alexa.com").Get("/").Reply(200).BodyString("<h1>API</h1>")
	gock.New("https://www.alexa.com").Get("/").Reply(200).BodyString("<h1>API</h1>")

	grabberAPIURLs := DefineExternalServices()

	for _, service := range grabberAPIURLs.Services {
		fmt.Printf("Grabbing %v\n", service.URL)
		// TODO: Implement Grab function for each External Service.
		// Need to figure out how to pass ExternalSeriveAPIURL{} object
		// to Grab function and satisfy Grabber interface.
	}
	str, err := Grab(d)
	if err != nil {
		fmt.Println("fail", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, 200, str)
}
