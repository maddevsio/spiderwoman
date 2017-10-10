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

type ExternalServiceList struct {
	Services []ExternalServiceItem
}

type Grabber interface {
	CheckConnection() (int, error)
	GetRawData() (string, error)
	//ParseData() (string, error)
}

type AlexaGrabber struct {
	ExternalServiceItem
}

type AhrefsGrabber struct {
	ExternalServiceItem
}

func GeneralCheckConnection(e ExternalServiceItem) (int, error) {
	resp, err := http.Get(e.URL)
	if err != nil {
		fmt.Println(err)
		return resp.StatusCode, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func GeneralGetRawData(e ExternalServiceItem) (string, error) {
	resp, err := http.Get(e.URL)
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

func (d AlexaGrabber) CheckConnection() (int, error) {
	return GeneralCheckConnection(d.ExternalServiceItem)
}

func (d AlexaGrabber) GetRawData() (string, error) {
	return GeneralGetRawData(d.ExternalServiceItem)
}

func (d AhrefsGrabber) CheckConnection() (int, error) {
	return GeneralCheckConnection(d.ExternalServiceItem)
}

func (d AhrefsGrabber) GetRawData() (string, error) {
	return GeneralGetRawData(d.ExternalServiceItem)
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

	//grabberAPIURLs := DefineExternalServices()

	alexaGrabber := AlexaGrabber{ExternalServiceItem{URL:"https://www.alexa.com", Name:"Alexa"}}
	ahrefsGrabber := AhrefsGrabber{ExternalServiceItem{URL:"https://www.ahrefs.com", Name:"Alexa"}}

	grabbers := []interface{}{&alexaGrabber, &ahrefsGrabber}

	for _, service := range grabbers {
		//fmt.Printf("Grabbing %v\n", service.(ExternalServiceItem).URL)
		str, err := Grab(service.(Grabber))
		if err != nil {
			fmt.Println("fail", err)
		}
		assert.NoError(t, err)
		assert.Equal(t, 200, str)
	}
}
