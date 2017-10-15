package main

import (
	"sync"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/maddevsio/simple-config"
	"github.com/urfave/cli"
	"os"
	"github.com/maddevsio/spiderwoman/lib"
)

type Ext struct {
	*gocrawl.DefaultExtender
}

var (
	mutex                 sync.Mutex
	hosts                 []string
	StopHosts             []lib.StopHostItem
	syncResolve           sync.WaitGroup
	err                   error
	externalLinksIterator int

	externalLinks         map[string]map[string]int
	externalLinksResolved map[string]map[string]int
	config 		      simple_config.SimpleConfig = simple_config.NewSimpleConfig("./config", "yml")

	userAgent             string                    = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
	resolveURLsPool       int                       = 100
	verbose               bool                      = true
	maxVisits             int                       = 10
	resolveTimeout        int                       = 30
	sqliteDBPath          string                    = config.GetString("db-path")
	excelFilePath         string                    = config.GetString("xls-path")
	internalOutPatterns   []string                  = []string{"/go/", "/go.php?", "/goto/", "/banners/click/", "/adrotate-out.php?", "/bsdb/bs.php?"}
	badSuffixes           []string                  = []string{".png", ".jpg", ".pdf"}
)

type Path struct {
	SqliteDBPath           string
	SourcesFilePath        string
	SourcesDefaultFilePath string
	TypesFilePath          string
	TypesDefaultFilePath   string
}

func main() {
	app := cli.NewApp()
	app.Name = "Spiderwoman"
	app.Usage = "Vertical crawler, which main target is to count links (resolved, e.g. from bit.ly) to external domains from all pages of given resources"

	app.Commands = []cli.Command{
		{
			Name:    "once",
			Aliases: []string{"o"},
			Usage:   "run the crawl and stop",
			Action:  actionOnce,
		},
		{
			Name:    "forever",
			Aliases: []string{"f"},
			Usage:   "start crawl forever using cron feature",
			Action:  actionForever,
		},
		{
			Name:    "excel",
			Aliases: []string{"e"},
			Usage:   "only create xls file",
			Action:  actionExcel,
		},
		{
			Name:    "grab",
			Aliases: []string{"g"},
			Usage:   "use grabber service only",
			Action:  actionGrab,
		},
	}

	app.Run(os.Args)
}