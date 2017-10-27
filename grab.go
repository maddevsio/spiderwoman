package main

import (
	"fmt"
	"log"

	"github.com/maddevsio/spiderwoman/grabber"
	"github.com/maddevsio/spiderwoman/lib"
)

func grab(path Path) {
	featuredHosts, err := lib.GetFeaturedHosts(path.SqliteDBPath)
	if err != nil {
		fmt.Println("Error gettings featured hosts", err)
	}

	alexaGrabber := grabber.AlexaGrabber{grabber.Service{Name: "Alexa"}}
	whoisGrabber := grabber.WhoisGrabber{grabber.Service{Name: "Whois"}}

	grabbers := []interface{}{&alexaGrabber, &whoisGrabber}
	for _, host := range featuredHosts {
		fmt.Printf("\nGrabbing %v data\n", host)
		for _, service := range grabbers {
			_, err := grabber.GrabAndSave(service.(grabber.Grabber), host, path.SqliteDBPath)
			if err != nil {
				log.Fatalf("GRABBER error %v", err)
			}
		}
	}
}
