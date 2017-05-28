package main

import (
	"log"

	"github.com/gin-contrib/gzip"
	"github.com/maddevsio/simple-config"
	"github.com/maddevsio/spiderwoman/lib"
	"gopkg.in/gin-gonic/gin.v1"
)

func GetAPIEngine(config simple_config.SimpleConfig) *gin.Engine {
	lib.CreateDBIfNotExists(config.GetString("db-path"))
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.BestCompression))
	accounts := gin.Accounts{config.GetString("admin-user"): config.GetString("admin-password")}

	r.Use(gin.BasicAuth(accounts))
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
	r.Static("/images", "./images")
	r.Static("/xls", config.GetString("xls-dir"))

	// main page
	r.GET("/", func(c *gin.Context) {
		dates, _ := lib.GetAllDaysFromMonitor(config.GetString("db-path"))
		s, _ := lib.GetCrawlStatus(config.GetString("db-path"))
		c.HTML(200, "index.html", gin.H{
			"title":  "Spiderwoman",
			"status": s,
			"dates":  dates,
			"dateQS": c.Query("date"), // pass this param to the "index.html" template
			"newQS":  c.Query("new"),
		})
	})

	// get json with monitor data to show in html table
	r.GET("/all", func(c *gin.Context) {
		var m []lib.Monitor
		if c.Query("date") != "" {
			if c.Query("new") == "1" {
				// this call is for new findings by a given date
				// e.g /?date=2006-12-04&new=1
				m, _ = lib.GetNewExtractedHostsForDay(config.GetString("db-path"), c.Query("date"))
			} else {
				// if we do not have "new" parameter, that just get all crawled
				// data for a given date
				m, _ = lib.GetAllDataFromMonitorByDay(config.GetString("db-path"), c.Query("date"))
			}
		} else {
			// in this case, when we do not have "date" parameter
			// we return all data from DB
			// TODO: for now this is useless and need to be deleted
			// and removed from the web interface
			m, _ = lib.GetAllDataFromMonitor(config.GetString("db-path"), 9)
		}
		c.JSON(200, m)
	})

	r.GET("/all-for-host", func(c *gin.Context) {
		var m []lib.Monitor
		if c.Query("host") != "" {
			m, _ = lib.GetAllDataFromMonitorByExternalHost(config.GetString("db-path"), c.Query("host"))
		}
		c.JSON(200, m)
	})

	// this is test endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r
}

func main() {
	// read the  config and run gin
	config := simple_config.NewSimpleConfig("../config", "yml")
	log.Printf("Server started on %v", config.GetString("api-port"))
	GetAPIEngine(config).Run(config.GetString("api-port"))
}
