package main

import "gopkg.in/gin-gonic/gin.v1"

func GetAPIEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return r
}

func main() {
	GetAPIEngine().Run(":8080")
}
