package main

import (
	"log"
	"sync"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/maddevsio/spiderwoman/lib"
	"os"
	"strings"
)

func crawl(path Path) {
	StopHosts = nil // clear this slice on every crawl
	StopHosts, err = lib.GetStopHosts(path.SqliteDBPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	externalLinks = make(map[string]map[string]int)
	externalLinksResolved = make(map[string]map[string]int)
	lib.SetCrawlStatus(path.SqliteDBPath, "Crawl started and crawling")
	hosts, err = lib.GetHostsFromFile(path.SourcesFilePath, path.SourcesDefaultFilePath)
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

	lib.SetCrawlStatus(path.SqliteDBPath, "Resolving URLS")
	log.Print("Going to resolve URLs...")
	for host := range externalLinks {
		for url, times := range externalLinks[host] {
			externalLinksIterator++
			syncResolve.Add(1)
			go func(url string, times int, host string, wg *sync.WaitGroup, mutex *sync.Mutex) {
				resolvedUrl := strings.ToLower(lib.Resolve(url, host, resolveTimeout, verbose, userAgent, mutex))
				defer wg.Done()
				if lib.HasStopHost(resolvedUrl, StopHosts) {
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

	lib.SetCrawlStatus(path.SqliteDBPath, "Saving the list")
	log.Print("Saving the list")
	lib.SaveDataToSqlite(path.SqliteDBPath, externalLinksResolved, verbose)
	lib.SetCrawlStatus(path.SqliteDBPath, "Crawl done")

	// no worries, we can rewrite all xls files on every call, this is not critical
	// when we will have more than 100 days of data, than we can think about optimization
	// createAllXLSByDays(path.SqliteDBPath)
}
