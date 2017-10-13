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

const (
	dbName = "spiderwoman-test"
)

type Service struct {
	Name string
	URL  string
}

type Grabber interface {
	CheckConnection() (int, error)
	GetServiceInfo() Service
	Do(featuredHost string) (string, error)
}

type AlexaGrabber struct {
	Service
}

type AhrefsGrabber struct {
	Service
}

func GeneralCheckConnection(s Service) (int, error) {
	resp, err := http.Get(s.URL)
	if err != nil {
		fmt.Println(err)
		return resp.StatusCode, err
	}
	defer resp.Body.Close()
	return 200, nil
}

func GeneralDo(s Service) (string, error) {
	resp, err := http.Get(s.URL)
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

func (ag AlexaGrabber) CheckConnection() (int, error) {
	// No need to check connection, because we use library for Alexa
	return 200, nil
}

func (ag AlexaGrabber) GetServiceInfo() Service {
	return ag.Service
}

func (ag AlexaGrabber) Do(featuredHost string) (string, error) {
	globalRank, err := alexa.GlobalRank(featuredHost)
	if err != nil {
		fmt.Printf("Alexa.Do(): %s\n", err)
		return "", err
	}
	fmt.Printf("Alexa.Do(): %s rank in alexa is %s\n", featuredHost, globalRank)
	return globalRank, nil
}

func (hg AhrefsGrabber) CheckConnection() (int, error) {
	return GeneralCheckConnection(hg.Service)
}

func (hg AhrefsGrabber) Do(featuredHost string) (string, error) {
	return GeneralDo(hg.Service)
}

func (hg AhrefsGrabber) GetServiceInfo() Service {
	return hg.Service
}

func GrabAndSave(g Grabber, featuredHost string) (bool, error) {
	conn, err := g.CheckConnection()
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	if conn == 200 {
		rawData, err := g.Do(featuredHost)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		fmt.Println(rawData)
		gd := lib.GrabberData{}

		gd.Service = g.GetServiceInfo().Name
		gd.Host = featuredHost
		gd.Data = rawData
		success := lib.SaveGrabbedData(dbName, gd)
		return success, nil
	}
	return false, nil
}

func TestGrabber(t *testing.T) {
	defer gock.Off()
	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExists(dbName)
	_, err := lib.AddFeaturedHost(dbName, "namba.kg")
	_, err = lib.AddFeaturedHost(dbName, "ts.kg")
	_, err = lib.AddFeaturedHost(dbName, "diesel.elcat.kg")
	featuredHosts, err := lib.GetFeaturedHosts(dbName)
	if err != nil {
		fmt.Println("Error gettings featured hosts", err)
	}

	alexaGrabber := AlexaGrabber{Service{Name: "Alexa"}}
	// ahrefsGrabber := AhrefsGrabber{Service{URL: "https://www.ahrefs.com", Name: "Alexa"}}

	grabbers := []interface{}{&alexaGrabber}
	for _, host := range featuredHosts {
		fmt.Printf("\nGrabbing %v data\n", host)
		for _, service := range grabbers {
			gock.New("http://data.alexa.com").Get("/").Reply(200).BodyString(fmt.Sprintf("<POPULARITY URL='%s' TEXT='12345' SOURCE='panel'/>", host))
			success, err := GrabAndSave(service.(Grabber), host)
			if err != nil {
				fmt.Println("GrabAndSave(): ", err)
			}
			assert.NoError(t, err)
			assert.Equal(t, true, success)
		}
	}
}
