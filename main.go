package main

import (
	"strings"
	"sync"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/maddevsio/spiderwoman/lib"
	"github.com/jasonlvhit/gocron"
	"github.com/maddevsio/simple-config"
	"log"
	"github.com/urfave/cli"
	"os"
)

type Ext struct {
	*gocrawl.DefaultExtender
}

var (
	mutex                 sync.Mutex
	hosts                 []string
	stopHosts             []string
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
	excelZipFilePath      string                    = config.GetString("zip-xls-path")
	internalOutPatterns   []string                  = []string{"/go/", "/go.php?", "/goto/", "/banners/click/", "/adrotate-out.php?", "/bsdb/bs.php?"}
	badSuffixes           []string                  = []string{".png", ".jpg", ".pdf"}
)

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
	}

	app.Run(os.Args)
}

func actionOnce(c *cli.Context) error {
	initialize()
	crawl()
	return nil
}

func actionForever(c *cli.Context) error {
	initialize()
	log.Print("All is OK. Starting cron job...")
	if config.GetString("box") == "dev" {
		log.Print("This is a dev box")
		gocron.Every(1).Minute().Do(crawl) // this is for testing on dev box
	} else {
		log.Print("This is production")
		if config.GetString("start-time") == "" {
			log.Fatal("You need to set start-time value in config.yaml")
		}
		gocron.Every(1).Day().At(config.GetString("start-time")).Do(crawl)
	}
	<- gocron.Start()
	return nil
}

func actionExcel(c *cli.Context) error {
	initialize()
	createXLS_BackupDB_Zip()
	return nil
}

func initialize() {
	lib.CreateDBIfNotExists(sqliteDBPath)
	lib.ClearResolveCache()
	err = lib.PopulateHostsAndTypes(sqliteDBPath, lib.SitesFilepath, lib.SitesDefaultFilepath)
	if err != nil {
		log.Fatal("Types population error")
	}
}

func crawl() {
	externalLinks = make(map[string]map[string]int)
	externalLinksResolved = make(map[string]map[string]int)
	lib.SetCrawlStatus(sqliteDBPath, "Crawl started and crawling")
	hosts, err = lib.GetHostsFromFile(lib.SitesFilepath, lib.SitesDefaultFilepath)
	if err != nil {
		log.Printf("Error opening or parsing config file: %v", err)
		return
	}

	for _, host := range hosts {
		ext := &Ext{&gocrawl.DefaultExtender{}}
		opts := gocrawl.NewOptions(ext)
		opts.CrawlDelay = 0
		if verbose {
			opts.LogFlags = gocrawl.LogAll
		} else {
			opts.LogFlags = gocrawl.LogError
		}
		opts.SameHostOnly = true
		opts.MaxVisits = maxVisits
		opts.HeadBeforeGet = false
		opts.UserAgent = userAgent
		opts.RobotUserAgent = userAgent
		c := gocrawl.NewCrawlerWithOptions(opts)
		c.Run("http://" + host)
	}

	lib.SetCrawlStatus(sqliteDBPath, "Resolving URLS")
	log.Print("Going to resolve URLs...")
	for host := range externalLinks {
		for url, times := range externalLinks[host] {
			externalLinksIterator++
			syncResolve.Add(1)
			go func(url string, times int, host string, wg *sync.WaitGroup, mutex *sync.Mutex) {
				resolvedUrl := lib.Resolve(url, host, resolveTimeout, verbose, userAgent, mutex)
				defer wg.Done()

				if lib.HasStopHost(resolvedUrl, stopHosts) {
					log.Printf("Url %v is in stoplist, not saving in map", resolvedUrl)
					return
				}

				mutex.Lock()
				if externalLinksResolved[host] == nil {
					externalLinksResolved[host] = make(map[string]int)
				}
				externalLinksResolved[host][resolvedUrl] = times
				mutex.Unlock()
			}(url, times, host, &syncResolve, &mutex)
			if externalLinksIterator%resolveURLsPool == 0 {
				syncResolve.Wait()
			}
		}
	}
	syncResolve.Wait()

	lib.SetCrawlStatus(sqliteDBPath, "Saving the list")
	log.Print("Saving the list")
	lib.SaveDataToSqlite(sqliteDBPath, externalLinksResolved, verbose)
	lib.SetCrawlStatus(sqliteDBPath, "Crawl done")

	createXLS_BackupDB_Zip()
}

func createXLS_BackupDB_Zip() {
	days, _ := lib.GetAllDaysFromMonitor(sqliteDBPath)
	log.Printf("Appendig XLS file with sheet %v", days[0])
	err = lib.AppendExcelFromDB(sqliteDBPath, excelFilePath, days[0])
	if (err != nil && strings.Contains(err.Error(), "no such file or directory")) {
		lib.CreateEmptyExcel(excelFilePath)
		log.Print("Trying to create all sheets in excel file")
		for _, day := range days {
			log.Printf("Appendig XLS file with sheet %v", day)
			err = lib.AppendExcelFromDB(sqliteDBPath, excelFilePath, day)
			if err != nil {
				log.Print(err)
			}
		}
	}

	log.Print("Backuping database")
	err = lib.BackupDatabase(sqliteDBPath)
	if (err != nil) {
		log.Printf("Backup error: %v", err)
	} else {
		log.Print("Database has been copied to /tmp/res.db")
	}

	log.Print("Zip XLS File")
	err = lib.ZipFile(excelFilePath, excelZipFilePath)
	if (err != nil) {
		log.Printf("Zip error: %v", err)
	} else {
		log.Printf("Zipped xls file was saved in %v", excelZipFilePath)
	}
}
