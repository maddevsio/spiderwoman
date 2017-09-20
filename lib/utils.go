package lib

import (
	"bufio"
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	SourcesFilePath        = "./sources.txt"
	SourcesDefaultFilePath = "./sources.default.txt"
	TypesFilePath          = "./types.txt"
	TypesHDefaultFilePath  = "./types.default.txt"
	StopsFilePath          = "./stops.txt"
	StopsDefaultFilePath   = "./stops.default.txt"
)

var (
	resolveCache     map[string]string
	lastCachedReturn = false
)

func ClearResolveCache() {
	resolveCache = make(map[string]string)
}

func Debug(data []byte, err error) {
	if err == nil {
		log.Printf("%s\n\n", data)
	} else {
		log.Printf("%s\n\n", err)
	}
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

				re := regexp.MustCompile("^www\\.")
				externalHost = re.ReplaceAllString(externalHost, "")

				re = regexp.MustCompile(":443$")
				externalHost = re.ReplaceAllString(externalHost, "")

				if u.Host == "" {
					continue
				}
			}
			if verbose {
				log.Printf("Saving result of %s: ", externalLink)
			}
			m := Monitor{}
			m.SourceHost = sourceHost
			m.ExternalHost = externalHost
			m.Count = count
			m.ExternalLink = externalLink
			res := SaveRecordToMonitor(DBFilepath, m)
			if verbose {
				log.Printf("The result of saving is: %t", res)
			}
		}
	}
	return true
}

// TODO: need to use cache, do not resolve same URLs
func Resolve(url string, host string, resolveTimeout int, verbose bool, userAgent string, mutex *sync.Mutex) string {
	lastCachedReturn = false
	if resolveCache[url] != "" {
		log.Printf("URL %v is in cache, return the resolved value %v", url, resolveCache[url])
		lastCachedReturn = true
		return resolveCache[url]
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(resolveTimeout) * time.Second,
	}

	if verbose {
		log.Println("Initial URL " + url)
	}

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		if verbose {
			log.Println("Bad URL: " + url + " Err:" + err.Error())
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
			dump, errDump := httputil.DumpResponse(response, false)
			if errDump == nil {
				Debug(dump, nil)
			}
			log.Printf("Resolved URL %v", response.Request.URL.String())
		}
		defer response.Body.Close()

		if mutex != nil {
			mutex.Lock()
		}

		resolveCache[url] = strings.ToLower(response.Request.URL.String())

		if mutex != nil {
			mutex.Unlock()
		}

		return strings.ToLower(response.Request.URL.String())
	} else {
		log.Printf("Error client.Do %v", err)
		return url
	}
}

func GetHostsFromFile(sourcesFilePath string, sourcesDefaultFilePath string) ([]string, error) {
	hosts, err := GetSliceFromFile(sourcesFilePath, sourcesDefaultFilePath)
	if err != nil {
		return []string{}, err
	}
	return hosts, nil
}

func HasStopHost(DBFilePath string, href string) bool {
	// TODO: need to get the list of stops once per crawl, not per check
	// stopHosts, _ = GetSliceFromFile(StopsFilePath, StopsDefaultFilePath)
	stopHosts, err := GetStopHosts(DBFilePath)
	if err != nil {
		log.Println(err)
		return false
	}

	for _, hostItem := range stopHosts {
		if strings.Contains(strings.ToLower(href), strings.ToLower(hostItem.Host)) {
			return true
		}
	}

	// for i := range stopHosts {
	// 	if strings.Contains(strings.ToLower(href), strings.ToLower(stopHosts[i])) {
	// 		return true
	// 	}
	// }
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

func BackupDatabase(dbPath string) error {
	srcFolder := dbPath
	destFolder := "/tmp/res.db"
	cpCmd := exec.Command("cp", "-r", srcFolder, destFolder)
	return cpCmd.Run()
}

func ZipFile(excelFilePath string, zipFilePath string) error {
	cpCmd := exec.Command("zip", zipFilePath, excelFilePath)
	return cpCmd.Run()
}

func PopulateHostsAndTypes(DBFilePath string, typesFilePath string, typesDefaultFilePath string) error {
	lines, err := GetSliceFromFile(typesFilePath, typesDefaultFilePath)
	if err != nil {
		log.Print(err)
		return err
	}
	err = DeleteTypesTable(DBFilePath)
	if err != nil {
		log.Print(err)
		return err
	}
	for _, line := range lines {
		hostName := strings.TrimSpace(strings.Split(line, " ")[0])
		hostType := strings.TrimSpace(strings.Split(line, " ")[1])
		err := SaveHostType(DBFilePath, hostName, hostType)
		if err != nil {
			log.Print(err)
			return err
		} else {
			log.Printf("Host %v and type %v saved", hostName, hostType)
		}
	}
	return nil
}

func MigrateStopHosts(DBName string, stopsFilePath string) error {
	file, err := os.Open(stopsFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err := AddStopHost(DBName, scanner.Text())
		if err != nil {
			log.Println(err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return nil
}
