package lib

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	SitesFilepath        = "./sites.txt"
	SitesDefaultFilepath = "./sites.default.txt"
	StopsFilepath        = "./stops.txt"
	StopsDefaultFilepath = "./stops.default.txt"
)

func Debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		fmt.Printf("%s\n\n", err)
	}
}

func SaveFile(data []string, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("%s %s\n", "Cannot create the file", err)
		return
	}
	defer file.Close()
	for _, line := range data {
		fmt.Fprintf(file, line)
	}
	fmt.Printf("File %s created", filename)
}

func GetSliceFromFile(realFile string, defaultFile string) ([]string, error) {
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

func SaveDataToSqlite(DBFilepath string, externalLinksResolved map[string]map[string]int, verbose bool) bool {
	for sourceHost, externalLinks := range externalLinksResolved {
		for externalLink, count := range externalLinks {
			var externalHost string
			u, err := url.Parse(externalLink)
			if err != nil {
				externalHost = externalLink
			} else {
				externalHost = u.Host
			}
			if verbose {
				fmt.Printf("Saving result of %s is: ", externalLink)
			}
			res := SaveRecordToMonitor(DBFilepath, sourceHost, externalLink, count, externalHost)
			if verbose {
				fmt.Printf("%t\n", res)
			}
		}
	}
	return true
}

// TODO: need to use cache, do not resolve same URLs
func Resolve(url string, host string, resolveTimeout int, verbose bool, userAgent string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(resolveTimeout) * time.Second,
	}

	if verbose {
		fmt.Println("Initial URL " + url)
	}

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		if verbose {
			fmt.Println("Bad URL: " + url + " Err:" + err.Error())
		}
		return url
	}

	request.Header.Add("User-Agent", userAgent)
	request.Header.Add("Referer", "http://"+host)

	if verbose {
		dump, errDump := httputil.DumpRequestOut(request, false)
		if errDump == nil {
			Debug(dump, nil)
		}
	}

	response, err := client.Do(request)
	if err == nil {
		if verbose {
			fmt.Println("Resolved URL " + response.Request.URL.String())
		}
		defer response.Body.Close()
		return response.Request.URL.String()
	} else {
		fmt.Println("Error client.Do" + err.Error())
		return url
	}
}

func GetHostsFromFile() ([]string, error) {
	return GetSliceFromFile(SitesFilepath, SitesDefaultFilepath)
}

func HasStopHost(href string, stopHosts []string) bool {
	if len(stopHosts) == 0 {
		stopHosts, _ = GetSliceFromFile(StopsFilepath, StopsDefaultFilepath)
	}

	for i := range stopHosts {
		if strings.Contains(strings.ToLower(href), strings.ToLower(stopHosts[i])) {
			return true
		}
	}
	return false
}

func HasInternalOutPatterns(href string, internalOutPatterns []string) bool {
	for i := range internalOutPatterns {
		if strings.Contains(href, internalOutPatterns[i]) {
			return true
		}
	}
	return false
}

func HasBadSuffixes(href string, badSuffixes []string) bool {
	for i := range badSuffixes {
		if strings.HasSuffix(href, badSuffixes[i]) {
			return true
		}
	}
	return false
}
