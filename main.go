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
	"os"
	"bufio"
	"net/http/httputil"
	"log"
	"net/url"
)

var (
	mutex sync.Mutex
	hosts []string
	stopHosts []string
	syncCrawl sync.WaitGroup
	syncResolve sync.WaitGroup
	err error
	externalLinksIterator int

	externalLinks map[string]map[string]int         = make(map[string]map[string]int)
	externalLinksResolved map[string]map[string]int = make(map[string]map[string]int)
	userAgent string = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
	verbose bool     = false
	maxVisits int    = 30
	internalOutPatterns []string = []string{"/go/", "/go.php?", "/goto/", "/banners/click/", "/adrotate-out.php?", "/bsdb/bs.php?"}
	badSuffixes []string = []string{".png", ".jpg", ".pdf"}
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

func main() {
	hosts, err = getHostsFromFile()
	if err != nil {
		fmt.Println("Error opening or parsing config file: " + err.Error())
		return
	}

	for i := 0; i < len(hosts); i++ {
		syncCrawl.Add(1)
		go func(key int) {
			fmt.Println(hosts[key])
			ext := &Ext{&gocrawl.DefaultExtender{}}
			// Set custom options
			opts := gocrawl.NewOptions(ext)
			opts.CrawlDelay = 0
			opts.LogFlags = gocrawl.LogError
			opts.SameHostOnly = true
			opts.MaxVisits = maxVisits
			opts.UserAgent = userAgent
			opts.RobotUserAgent = userAgent
			c := gocrawl.NewCrawlerWithOptions(opts)
			defer syncCrawl.Done()
			c.Run(hosts[key])
		}(i)
	}
	syncCrawl.Wait()

	//spew.Dump(externalLinks)

	fmt.Println("Going to resolve URLs...")
	for host := range externalLinks {
		for url, times := range externalLinks[host] {
			externalLinksIterator++
			syncResolve.Add(1)
			go func(url string, times int, host string, wg *sync.WaitGroup) {
				resolvedUrl := resolve(url, host)

				mutex.Lock()
				if (externalLinksResolved[host] == nil) {
					externalLinksResolved[host]= make(map[string]int)
				}
				externalLinksResolved[host][resolvedUrl] = times
				mutex.Unlock()

				wg.Done()
			}(url, times, host, &syncResolve)
			if externalLinksIterator%10 == 0 {
				syncResolve.Wait()
			}
		}
	}
	syncResolve.Wait()

	//spew.Dump(externalLinksResolved)

	fmt.Println("\n\n\n\nSorting the list")
	sortMapByKeys(externalLinksResolved)
}

// helper functions

func sortMapByKeys(externalLinksResolved map[string]map[string]int) {
	for host, m := range externalLinksResolved {
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
				var externalLinkHost string
				u, err := url.Parse(s)
				if err !=nil {
					externalLinkHost = s
				} else {
					externalLinkHost = u.Host
				}
				fmt.Printf("%s\t%s\t%d\t%s\n", host, s, k, externalLinkHost)
			}
		}
	}
}

// TODO: need to use cache, do not resolve same URLs
func resolve(url string, host string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 15 * time.Second}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		if verbose {
			fmt.Println("Bad URL: " + url + " Err:" + err.Error())
		}
		return url
	}
	request.Header.Add("User-Agent", userAgent)
	request.Header.Add("Referer", "http://" + host)

	if verbose {
		debug(httputil.DumpRequestOut(request, false))
	}

	response, err := client.Do(request)
	if err == nil {
		fmt.Println(response.Request.URL.String())
		return response.Request.URL.String()
	} else {
		return url
	}
}

func getHostsFromFile() ([]string, error) {
	return getSliceFromFile("./sites.txt", "./sites.default.txt")
}

func hasStopHost(href string) bool {
	if len(stopHosts) == 0 {
		stopHosts, err = getSliceFromFile("./stops.txt", "./stops.default.txt")
	}

	for i := range stopHosts {
		if (strings.Contains(href, stopHosts[i])) {
			return true
		}
	}
	return false
}

func hasInternalOutPatterns(href string) bool {
	for i := range internalOutPatterns {
		if strings.Contains(href, internalOutPatterns[i]) {
			return true;
		}
	}
	return false;
}

func hasBadSuffixes(href string) bool {
	for i := range badSuffixes {
		if strings.HasSuffix(href, badSuffixes[i]) {
			return true;
		}
	}
	return false;
}

func getSliceFromFile(realFile string, defaultFile string) ([]string, error) {
	file, err := os.Open(realFile)
	if err != nil {
		file, err = os.Open(defaultFile)
		if err != nil {
			return nil, err
		}
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}