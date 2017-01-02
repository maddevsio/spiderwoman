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
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/all", func(c *gin.Context) {
		m, _ := lib.GetAllDataFromMonitor(dbPath)
		c.JSON(200, m)
	})

	return r
}

func main() {
	GetAPIEngine("../res.db").Run(":8080")
	// TODO: extract res.db in var or const and use in from lib package
}
