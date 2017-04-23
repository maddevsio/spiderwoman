package main

import (
	"strings"
	"time"
	"github.com/PuerkitoBio/goquery"
	"github.com/PuerkitoBio/gocrawl"
	"github.com/maddevsio/spiderwoman/lib"
	"net/http"
	"log"
)

func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	log.Printf("Visit: %s\n", ctx.URL())
	if doc == nil {
		return nil, true
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		href = strings.ToLower(href)

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