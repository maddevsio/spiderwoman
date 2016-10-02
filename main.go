package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"time"
	"strings"
	"crypto/tls"
	"sync"
	"sort"
)

type Ext struct {
	*gocrawl.DefaultExtender
}

var (
	m map[string]int
	mResolved map[string]int
	mutex sync.Mutex
)

func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	fmt.Printf("Visit: %s\n", ctx.URL())
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")

		// restriction chain, shitcoded, will be refactored
		if strings.Contains(href, ctx.URL().Host) {
			return
		}
		if (!strings.HasPrefix(href, "http")) {
			return
		}
		if (strings.Contains(href, "//telegram.me")) {
			return
		}
		if (strings.Contains(href, "//plus.google.com")) {
			return
		}
		if (strings.Contains(href, "//www.facebook.com")) {
			return
		}
		if (strings.Contains(href, "//vk.com")) {
			return
		}
		if (strings.Contains(href, "//www.youtube.com")) {
			return
		}
		if (strings.Contains(href, "//twitter.com")) {
			return
		}
		if (strings.HasSuffix(href, ".png")) {
			return
		}
		if (strings.HasSuffix(href, ".jpg")) {
			return
		}

		m[href] += 1

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

func main() {
	var syncCrawl sync.WaitGroup
	var syncResolve sync.WaitGroup
	m = make(map[string]int)
	mResolved = make(map[string]int)
	var hosts [1]string
	hosts[0] = "http://nambataxi.kg/"
	for i:=0; i<len(hosts); i++ {
		syncCrawl.Add(1)
		go func(key int) {
			fmt.Println(hosts[key])
			ext := &Ext{&gocrawl.DefaultExtender{}}
			// Set custom options
			opts := gocrawl.NewOptions(ext)
			opts.CrawlDelay = 0
			opts.LogFlags = gocrawl.LogError
			opts.SameHostOnly = true
			opts.MaxVisits = 10
			opts.UserAgent = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
			opts.RobotUserAgent = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
			c := gocrawl.NewCrawlerWithOptions(opts)
			defer syncCrawl.Done()
			c.Run(hosts[key])
		}(i)
	}
	syncCrawl.Wait()

	fmt.Println("Going to resolve URLs...")
	for url, times := range m {
		syncResolve.Add(1)
		go func(url string, times int, wg *sync.WaitGroup) {
			resolvedUrl := resolve(url)

			mutex.Lock()
			mResolved[resolvedUrl] = times
			mutex.Unlock()

			wg.Done()
		}(url, times, &syncResolve)
	}
	syncResolve.Wait()

	fmt.Println("Sorting the list")
	sortMapByKeys(mResolved)
}

func sortMapByKeys(m map[string]int) {
	n := map[int][]string{}
	var a []int
	for k, v := range m {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	for _, k := range a {
		for _, s := range n[k] {
			fmt.Printf("%s, %d\n", s, k)
		}
	}
}

func resolve(url string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	response, err := client.Get(url)
	if err == nil {
		fmt.Println(response.Request.URL.String())
		return response.Request.URL.String()
	} else {
		return url
	}
}