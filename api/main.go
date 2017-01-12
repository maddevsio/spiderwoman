package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/maddevsio/spiderwoman/lib"
	"github.com/gin-contrib/gzip"
)

func GetAPIEngine(dbPath string) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.BestCompression))
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
	r.Static("/images", "./images")
	r.StaticFile("/spiderwoman.xls", "/tmp/spiderwoman.xls")

	r.GET("/", func(c *gin.Context) {
		s, _ := lib.GetCrawlStatus("../res.db")
		c.HTML(200, "index.html", gin.H{
			"title": "Spiderwoman",
			"status": s, // TODO extract res.db to the config!!11
		})
	})

	r.GET("/all", func(c *gin.Context) {
		m, _ := lib.GetAllDataFromMonitor(dbPath)
		c.JSON(200, m)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r
}

func main() {
	GetAPIEngine("../res.db").Run(":8080")
	// TODO: extract res.db and port to config
}
