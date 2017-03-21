package main

import (
	"github.com/urfave/cli"
	"github.com/jasonlvhit/gocron"
	"log"
	"github.com/maddevsio/spiderwoman/lib"
)

func actionOnce(c *cli.Context) error {
	path := Path{sqliteDBPath, lib.SitesFilepath, lib.SitesDefaultFilepath}
	initialize(path)
	crawl(path)
	return nil
}

func actionForever(c *cli.Context) error {
	path := Path{sqliteDBPath, lib.SitesFilepath, lib.SitesDefaultFilepath}
	initialize(path)
	log.Print("All is OK. Starting cron job...")
	if config.GetString("box") == "dev" {
		log.Print("This is a dev box")
		gocron.Every(1).Minute().Do(crawl) // this is for testing on dev box
	} else {
		log.Print("This is production")
		if config.GetString("start-time") == "" {
			log.Fatal("You need to set start-time value in config.yaml")
		}
		gocron.Every(1).Day().At(config.GetString("start-time")).Do(crawl, path)
	}
	<- gocron.Start()
	return nil
}

func actionExcel(c *cli.Context) error {
	path := Path{sqliteDBPath, lib.SitesFilepath, lib.SitesDefaultFilepath}
	initialize(path)
	createXLS_BackupDB_Zip()
	return nil
}

