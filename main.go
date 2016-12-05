package main

import (
	"fmt"
	"sync"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/maddevsio/spiderwoman/lib"
)

var (
	mutex sync.Mutex
	hosts []string
	stopHosts []string
	syncResolve sync.WaitGroup
	err error
	externalLinksIterator int

	externalLinks map[string]map[string]int         = make(map[string]map[string]int)
	externalLinksResolved map[string]map[string]int = make(map[string]map[string]int)
	userAgent string     = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
	resolveURLsPool int  = 100
	verbose bool         = true
	maxVisits int        = 10
	resolveTimeout int   = 30
	internalOutPatterns []string = []string{"/go/", "/go.php?", "/goto/", "/banners/click/", "/adrotate-out.php?", "/bsdb/bs.php?"}
	badSuffixes []string = []string{".png", ".jpg", ".pdf"}
)

func main() {
	hosts, err = lib.GetHostsFromFile()
	if err != nil {
		fmt.Println("Error opening or parsing config file: " + err.Error())
		return
	}

	for _, host := range hosts {
		fmt.Println(" hh " + host)
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
				if (externalLinksResolved[host] == nil) {
					externalLinksResolved[host]= make(map[string]int)
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

	fmt.Println("\n\n\n\nSorting the list")
	csv := lib.SortMapByKeys(externalLinksResolved, verbose)
	lib.SaveFile(csv, "res.csv")
}