package main

import (
	"github.com/urfave/cli"
	"github.com/jasonlvhit/gocron"
	"log"
	"github.com/maddevsio/spiderwoman/lib"
	"os"
)

var path = Path{sqliteDBPath, lib.SourcesFilePath, lib.SourcesDefaultFilePath, lib.TypesFilePath, lib.TypesHDefaultFilePath}

func actionOnce(c *cli.Context) error {
	initialize(path)
	crawl(path)
	return nil
}

func actionForever(c *cli.Context) error {
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
	err := os.MkdirAll(config.GetString("xls-dir"), 0777)
	if err != nil {
		log.Fatalf("cannot create dir for excel files: %v", err)
	}
	if c.Args().Get(0) != "noinit" {
		initialize(path)
	}
	createAllXLSByDays(path.SqliteDBPath)
	return nil
}

