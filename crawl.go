package main

import (
	"sync"
	"github.com/maddevsio/spiderwoman/lib"
	"log"
	"github.com/PuerkitoBio/gocrawl"
)

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

