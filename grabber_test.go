package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

type Grabber interface {
	CheckConnection() (int, error)
	GetRawData() (string, error)
	ParseData(string) (string, error)
}

type InputData struct {
	FeaturedHost          string
	ExternalServiceAPIUrl string
}

func (d InputData) CheckConnection() (int, error) {
	resp, err := http.Get(d.ExternalServiceAPIUrl)
	if err != nil {
		fmt.Println(err)
		return resp.StatusCode, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func (d InputData) GetRawData() (string, error) {
	query := fmt.Sprintf("%s/api/?host=%s", d.ExternalServiceAPIUrl, d.FeaturedHost)
	fmt.Println(query)
	resp, err := http.Get(query)
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

func (d InputData) ParseData(data string) (string, error) {
	return "here is parsed data", nil
}

func GrabAlexa(g Grabber) (int, error) {
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
		parsedData, err := g.ParseData(rawData)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(parsedData)
	}
	return conn, err
}

func TestGrabber(t *testing.T) {
	d := InputData{FeaturedHost: "nambataxi.kg", ExternalServiceAPIUrl: "https://www.alexa.com"}
	defer gock.Off()
	gock.New("https://www.alexa.com").Get("/").Reply(200).BodyString("<h1>API</h1>")
	gock.New("https://www.alexa.com").Get("/").Reply(200).BodyString("<h1>API</h1>")

	str, err := GrabAlexa(d)
	if err != nil {
		fmt.Println("fail", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, 200, str)
}
