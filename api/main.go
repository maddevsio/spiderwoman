package main

import (
	"html/template"
	"log"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/maddevsio/simple-config"
	"github.com/maddevsio/spiderwoman/lib"
)

func GetAPIEngine(config simple_config.SimpleConfig) *gin.Engine {
	lib.CreateDBIfNotExists(config.GetString("db-path"))

	// set the log level
	if config.GetString("box") != "dev" {
		log.Print("This is production")
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.Print("This is development")
		gin.SetMode(gin.DebugMode)
	}

	// Define templates
	templates := multitemplate.New()
	templates.AddFromFiles("index", "templates/base.html", "templates/index.html")
	templates.AddFromFiles("types", "templates/base.html", "templates/types.html")
	templates.AddFromFiles("report", "templates/base.html", "templates/report.html")
	templates.AddFromFiles("perfomance-report", "templates/base.html", "templates/perfomance_report.html")
	templates.AddFromFiles("stops", "templates/base.html", "templates/stops.html")

	r := gin.Default()
	r.HTMLRender = templates // Comment this if you want to use old templates
	// r.LoadHTMLGlob("templates/*") // Uncomment this if you want to use old templates
	r.Use(gzip.Gzip(gzip.BestCompression))
	accounts := gin.Accounts{config.GetString("admin-user"): config.GetString("admin-password")}

	r.Use(gin.BasicAuth(accounts))
	r.Static("/assets", "./assets")
	r.Static("/images", "./images")
	r.Static("/xls", config.GetString("xls-dir")) // this is for all existent XLS files

	// main page
	r.GET("/", func(c *gin.Context) {
		dates, _ := lib.GetAllDaysFromMonitor(config.GetString("db-path"), "7")
		s, _ := lib.GetCrawlStatus(config.GetString("db-path"))
		var types []string
		types, _ = lib.GetUniqueTypes(config.GetString("db-path"))
		featured_hosts, _ := lib.GetFeaturedHosts(config.GetString("db-path"))
		c.HTML(200, "index", gin.H{
			"title":          "Spiderwoman | Home",
			"status":         s,
			"dates":          dates,
			"dateQS":         c.Query("date"), // pass this param to the "index.html" template
			"newQS":          c.Query("new"),
			"types":          types,
			"featured_hosts": featured_hosts,
		})
	})

	// report page with metronic template
	r.GET("/report", func(c *gin.Context) {
		dates, _ := lib.GetAllDaysFromMonitor(config.GetString("db-path"), "")
		s, _ := lib.GetCrawlStatus(config.GetString("db-path"))
		var types []string
		types, _ = lib.GetUniqueTypes(config.GetString("db-path"))
		c.HTML(200, "report", gin.H{
			"title":     "Spiderwoman | Report",
			"status":    s,
			"dates":     dates,
			"endDate":   dates[0],
			"startDate": dates[len(dates)-1],
			"dateQS":    c.Query("date"), // pass this param to the "index.html" template
			"newQS":     c.Query("new"),
			"types":     types,
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
			// m, _ = lib.GetAllDataFromMonitor(config.GetString("db-path"), 9)
		}
		c.JSON(200, m)
	})

	// get all monitors filtered by host
	r.GET("/all-for-host", func(c *gin.Context) {
		var m []lib.Monitor
		if c.Query("host") != "" {
			m, _ = lib.GetAllDataFromMonitorByExternalHost(config.GetString("db-path"), c.Query("host"))
		}
		c.JSON(200, m)
	})

	// perfomance report page
	r.GET("/perfomance/", func(c *gin.Context) {
		s, _ := lib.GetCrawlStatus(config.GetString("db-path"))
		var data []lib.PerfomanceReportResponse
		var byTypes []lib.PerfomanceReportByHostTypeResponse
		var alexaRank []lib.GrabberData
		var whoisData lib.GrabberData
		if c.Query("host") != "" {
			data, _ = lib.PerfomanceReport(config.GetString("db-path"), c.Query("host"))
			byTypes, _ = lib.PerfomanceReportByHostTypes(config.GetString("db-path"), c.Query("host"))
			alexaRank, _ = lib.PerfomanceReportGrabberData(config.GetString("db-path"), "alexa", c.Query("host"))
			whoisData, _ = lib.PerfomanceReportLatestGrabberData(config.GetString("db-path"), "whois", c.Query("host"))
		}
		whoisDataHTML := template.HTML(strings.Replace(whoisData.Data, "\n", "<br>", -1))
		c.HTML(200, "perfomance-report", gin.H{
			"title":     "Spiderwoman | Perfomance Report For " + c.Query("host"),
			"status":    s,
			"host":      c.Query("host"),
			"data":      data,
			"byTypes":   byTypes,
			"alexaRank": alexaRank,
			"whoisDate": whoisData.Created,
			"whoisData": whoisDataHTML,
		})
	})

	// return xls on the fly for new data
	r.GET("/get-new-xls", func(c *gin.Context) {
		xlsFileName := "new-" + c.Query("date") + ".xls"
		xlsFilePath := "/tmp/" + xlsFileName
		lib.CreateExcelFromDB_NEW(config.GetString("db-path"), xlsFilePath, c.Query("date"))
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+xlsFileName)
		c.Header("Content-Type", "application/octet-stream")
		c.File(xlsFilePath)
	})

	// return xls on the fly for all data by the day
	r.GET("/get-day-xls", func(c *gin.Context) {
		xlsFileName := "day-" + c.Query("date") + ".xls"
		xlsFilePath := "/tmp/" + xlsFileName
		lib.CreateExcelFromDB(config.GetString("db-path"), xlsFilePath, c.Query("date"))
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+xlsFileName)
		c.Header("Content-Type", "application/octet-stream")
		c.File(xlsFilePath)
	})

	// this is test endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/types", func(c *gin.Context) {
		var t []lib.HostItem
		t, _ = lib.GetAllTypes(config.GetString("db-path"))
		c.HTML(200, "types", gin.H{
			"title": "Spiderwoman | Hosts Types",
			"hosts": t,
		})
	})

	r.GET("/types/delete", func(c *gin.Context) {
		err := lib.DeleteHost(config.GetString("db-path"), c.Query("host"))
		if err != nil {
			log.Fatal(err)
			c.JSON(500, nil)
		}
		c.JSON(200, nil)
	})

	r.POST("/types/create", func(c *gin.Context) {
		name := c.PostForm("host_name")
		host_type := c.PostForm("host_type")
		err := lib.SaveHostType(config.GetString("db-path"), name, host_type)
		if err != nil {
			log.Println(err)
			c.JSON(500, nil)
		}
		c.JSON(200, nil)
	})

	r.POST("/types/update/", func(c *gin.Context) {
		hostName := c.PostForm("host_name")
		hostType := c.PostForm("host_type")
		err := lib.UpdateOrCreateHostType(config.GetString("db-path"), hostName, hostType)
		if err != nil {
			log.Println(err)
			c.JSON(500, nil)
		}
		c.JSON(200, nil)
	})

	r.GET("/featured/add/", func(c *gin.Context) {
		msg, err := lib.AddFeaturedHost(config.GetString("db-path"), c.Query("host"))
		if err != nil {
			log.Println(err)
			c.JSON(500, nil)
		}
		if msg == "removed" {
			c.JSON(200, gin.H{
				"message": "Host is removed.",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "Host is added.",
			})
		}
	})

	r.GET("/featured/remove/", func(c *gin.Context) {
		err := lib.RemoveFeaturedHost(config.GetString("db-path"), c.Query("host"))
		if err != nil {
			log.Println(err)
			c.JSON(500, nil)
		}
		c.JSON(200, nil)
	})

	r.GET("/stops", func(c *gin.Context) {
		var st []lib.StopHostItem
		st, _ = lib.GetStopHosts(config.GetString("db-path"))
		c.HTML(200, "stops", gin.H{
			"title": "Spiderwoman | Stop Hosts",
			"hosts": st,
		})
	})

	r.GET("/stops/delete", func(c *gin.Context) {
		err := lib.RemoveStopHost(config.GetString("db-path"), c.Query("host"))
		if err != nil {
			log.Fatal(err)
			c.JSON(500, nil)
		}
		c.JSON(200, nil)
	})

	r.POST("/stops/create", func(c *gin.Context) {
		host := c.PostForm("host")
		err := lib.AddStopHost(config.GetString("db-path"), host)
		if err != nil {
			log.Println(err)
			c.JSON(500, nil)
		}
		c.JSON(200, nil)
	})
	return r
}

func main() {
	// read the  config and run gin
	config := simple_config.NewSimpleConfig("../config", "yml")
	log.Printf("Server started on %v", config.GetString("api-port"))
	GetAPIEngine(config).Run(config.GetString("api-port"))
}
