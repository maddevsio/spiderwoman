package main

import (
	"fmt"
	"strings"
	"time"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"github.com/PuerkitoBio/gocrawl"
)

type Ext struct {
	*gocrawl.DefaultExtender

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
			if !hasInternalOutPatterns(href) {
				return
			} else {
				if verbose {
					fmt.Println(href)
				}
			}

		}

		// analyze relative urls, e.g. /lolz.html
		if (!strings.HasPrefix(href, "http")) {
			if !hasInternalOutPatterns(href) {
				return
			} else {
				href = ctx.URL().Scheme + "://" + ctx.URL().Host + href
				if verbose {
					fmt.Println(href)
				}
			}
		}

		if (hasStopHost(href)) {
			return
		}

		if (hasBadSuffixes(href)) {
			return
		}

		mutex.Lock()
		if (externalLinks[ctx.URL().Host] == nil) {
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
