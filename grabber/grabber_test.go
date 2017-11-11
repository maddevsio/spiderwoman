package grabber

import (
	"testing"

	"github.com/maddevsio/spiderwoman/lib"
	"gopkg.in/h2non/gock.v1"

	"fmt"
	"github.com/stretchr/testify/assert"
)

const (
	dbName = "spiderwoman-test"
)

func TestGrabber(t *testing.T) {
	defer gock.Off()
	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExistsAndMigrate(dbName)
	_, err := lib.AddFeaturedHost(dbName, "namba.kg")
	_, err = lib.AddFeaturedHost(dbName, "ts.kg")
	_, err = lib.AddFeaturedHost(dbName, "diesel.elcat.kg")
	featuredHosts, err := lib.GetFeaturedHosts(dbName)
	if err != nil {
		fmt.Println("Error gettings featured hosts", err)
	}

	alexaGrabber := AlexaGrabber{Service{Name: "Alexa"}}

	grabbers := []interface{}{&alexaGrabber}
	for _, host := range featuredHosts {
		fmt.Printf("\nGrabbing %v data\n", host)
		for _, service := range grabbers {
			gock.New("http://data.alexa.com").Get("/").Reply(200).BodyString(fmt.Sprintf("<POPULARITY URL='%s' TEXT='12345' SOURCE='panel'/>", host))
			success, err := GrabAndSave(service.(Grabber), host, dbName)
			if err != nil {
				fmt.Println("GrabAndSave(): ", err)
			}
			assert.NoError(t, err)
			assert.Equal(t, true, success)
		}
	}
}
