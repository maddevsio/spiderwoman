package main

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"github.com/maddevsio/spiderwoman/lib"
	"github.com/jasonlvhit/gocron"
	"github.com/maddevsio/simple-config"
	"log"
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

	externalLinks         map[string]map[string]int  = make(map[string]map[string]int)
	externalLinksResolved map[string]map[string]int  = make(map[string]map[string]int)
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

func main() {
	log.Print("All is OK. Starting cron job...")
	if config.GetString("box") == "dev" {
		log.Print("This is a dev box")
		gocron.Every(1).Minute().Do(crawl) // this is for testing on dev box
	} else {
		log.Print("This is production")
		gocron.Every(1).Day().At("00:00").Do(crawl)
	}
	<- gocron.Start()
}

func crawl() {
	lib.CreateDBIfNotExists(sqliteDBPath)
	lib.SetCrawlStatus(sqliteDBPath, "Crawl started and crawling")
	hosts, err = lib.GetHostsFromFile()
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
		c.Run(host)
	}

	lib.SetCrawlStatus(sqliteDBPath, "Resolving URLS")
	log.Print("Going to resolve URLs...")
	for host := range externalLinks {
		for url, times := range externalLinks[host] {
			externalLinksIterator++
			syncResolve.Add(1)
			go func(url string, times int, host string, wg *sync.WaitGroup) {
				resolvedUrl := lib.Resolve(url, host, resolveTimeout, verbose, userAgent)

				mutex.Lock()
				if externalLinksResolved[host] == nil {
					externalLinksResolved[host] = make(map[string]int)
				}
				externalLinksResolved[host][resolvedUrl] = times
				mutex.Unlock()

				wg.Done()
			}(url, times, host, &syncResolve)
			if externalLinksIterator%resolveURLsPool == 0 {
				syncResolve.Wait()
			}
		}
	}
	syncResolve.Wait()

	lib.SetCrawlStatus(sqliteDBPath, "Saving the list")
	log.Print("Saving the list")
	lib.SaveDataToSqlite(sqliteDBPath, externalLinksResolved, verbose)
	lib.CreateExcelFromDB(sqliteDBPath, excelFilePath)
	lib.SetCrawlStatus(sqliteDBPath, "Crawl done")

	err := lib.BackupDatabase(sqliteDBPath)
	if (err != nil) {
		log.Printf("Backup error: %v", err)
	} else {
		log.Print("Database has been copied to /tmp/res.db")
	}
}

func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	log.Printf("Visit: %s\n", ctx.URL())
	if doc == nil {
		return nil, true
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")

		// analyze absolute urls, e.g. http://bla.com/lolz
		if strings.Contains(href, ctx.URL().Host) {
			if !lib.HasInternalOutPatterns(href, internalOutPatterns) {
				return
			} else {
				if verbose {
					log.Print(href)
				}
			}

		}

		// analyze relative urls, e.g. /lolz.html
		if !strings.HasPrefix(href, "http") {
			if !lib.HasInternalOutPatterns(href, internalOutPatterns) {
				return
			} else {
				href = ctx.URL().Scheme + "://" + ctx.URL().Host + href
				if verbose {
					log.Print(href)
				}
			}
		}

		if lib.HasStopHost(href, stopHosts) {
			return
		}

		if lib.HasBadSuffixes(href, badSuffixes) {
			return
		}

		mutex.Lock()
		if externalLinks[ctx.URL().Host] == nil {
			externalLinks[ctx.URL().Host] = make(map[string]int)
		}
		externalLinks[ctx.URL().Host][href] += 1
		mutex.Unlock()

	})
	return nil, true
}

func (e *Ext) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	return true
}

func (de *Ext) RequestRobots(ctx *gocrawl.URLContext, robotAgent string) (data []byte, doRequest bool) {
	return nil, false
}

func (e *Ext) ComputeDelay(host string, di *gocrawl.DelayInfo, lastFetch *gocrawl.FetchInfo) time.Duration {
	return 0
}
