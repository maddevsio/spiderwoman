package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/maddevsio/spiderwoman/lib"
	"github.com/negah/alexa"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

type ExternalServiceItem struct {
	Name string
	URL  string
}

type Grabber interface {
	CheckConnection() (int, error)
	Do(featuredHost string) (string, error)
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
	return 200, nil
}

func GeneralDo(e ExternalServiceItem) (string, error) {
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
	// No need to check connection, because we use library for Alexa
	return 200, nil
}

func (d AlexaGrabber) Do(featuredHost string) (string, error) {
	globalRank, err := alexa.GlobalRank(featuredHost)
	if err != nil {
		fmt.Printf("Do Alexa Grabber error: %s", err)
	} else {
		fmt.Printf("%s rank in alexa is %s\n", featuredHost, globalRank)
	}
	return globalRank, nil
}

func (d AhrefsGrabber) CheckConnection() (int, error) {
	return GeneralCheckConnection(d.ExternalServiceItem)
}

func (d AhrefsGrabber) Do(featuredHost string) (string, error) {
	return GeneralDo(d.ExternalServiceItem)
}

func GrabAndSave(g Grabber, featuredHost string) (int, error) {
	conn, err := g.CheckConnection()
	if err != nil {
		fmt.Println(err)
	}
	if conn == 200 {
		fmt.Println("Connection is ok, can proceed.")
		rawData, err := g.Do(featuredHost)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(rawData)
		// TODO: Save data returned from Do function to data base
	}
	return conn, err
}

func TestGrabber(t *testing.T) {
	defer gock.Off()
	dbName := "spiderwoman-test"
	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExists(dbName)
	err, _ := lib.AddFeaturedHost(dbName, "namba.kg")
	err, _ = lib.AddFeaturedHost(dbName, "ts.kg")
	err, _ = lib.AddFeaturedHost(dbName, "diesel.elcat.kg")
	featuredHosts, err := lib.GetFeaturedHosts(dbName)
	if err != nil {
		fmt.Println("Error gettings featured hosts", err)
	}

	alexaGrabber := AlexaGrabber{ExternalServiceItem{Name: "Alexa"}}
	// ahrefsGrabber := AhrefsGrabber{ExternalServiceItem{URL: "https://www.ahrefs.com", Name: "Alexa"}}

	grabbers := []interface{}{&alexaGrabber}
	for _, host := range featuredHosts {
		fmt.Printf("Grabbing %v data\n", host)
		for _, service := range grabbers {
			gock.New("http://data.alexa.com").Get("/").Reply(200).BodyString("<h1>API</h1>")
			str, err := GrabAndSave(service.(Grabber), host)
			if err != nil {
				fmt.Println("fail", err)
			}
			assert.NoError(t, err)
			assert.Equal(t, 200, str)
		}
	}
}
