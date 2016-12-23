package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"github.com/maddevsio/spiderwoman/lib"
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

	externalLinks         map[string]map[string]int = make(map[string]map[string]int)
	externalLinksResolved map[string]map[string]int = make(map[string]map[string]int)
	userAgent             string                    = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
	resolveURLsPool       int                       = 100
	verbose               bool                      = true
	maxVisits             int                       = 10
	resolveTimeout        int                       = 30
	sqliteDBPath          string                    = "./res.db"
	internalOutPatterns   []string                  = []string{"/go/", "/go.php?", "/goto/", "/banners/click/", "/adrotate-out.php?", "/bsdb/bs.php?"}
	badSuffixes           []string                  = []string{".png", ".jpg", ".pdf"}
)

func main() {
	hosts, err = lib.GetHostsFromFile()
	if err != nil {
		fmt.Println("Error opening or parsing config file: " + err.Error())
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

	fmt.Println("Going to resolve URLs...")
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

	fmt.Println("Saving the list")
	lib.CreateDBIfNotExists(sqliteDBPath)
	lib.SaveDataToSqlite(sqliteDBPath, externalLinksResolved, verbose)
}

func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	fmt.Printf("Visit: %s\n", ctx.URL())
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
					fmt.Println(href)
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
					fmt.Println(href)
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
