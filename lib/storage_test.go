package lib

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	DBFilepath = "spiderwoman-test"
)

func TestCreateDBIfNotExists(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
}

func TestCheckMonitorTable(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)

	db := getDB(DBFilepath)
	defer db.Close()

	rows, err := db.Query("SHOW TABLES LIKE 'monitor'")
	assert.NoError(t, err)
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		assert.Equal(t, nil, err)
		assert.Equal(t, "monitor", name)
	}

	err = rows.Err()
	assert.Equal(t, nil, err)
}

func TestCheckStatusTable(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)

	db := getDB(DBFilepath)
	defer db.Close()

	rows, err := db.Query("SHOW TABLES LIKE 'status'")
	assert.Equal(t, nil, err)
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		assert.Equal(t, nil, err)
		assert.Equal(t, "status", name)
	}

	err = rows.Err()
	assert.Equal(t, nil, err)
}

func TestSaveRecordToMonitor(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)

	m := Monitor{}
	m.SourceHost = "http://a"
	m.ExternalLink = "http://b/1"
	m.Count = 800
	m.ExternalHost = "b"
	res := SaveRecordToMonitor(DBFilepath, m)
	assert.Equal(t, true, res)

	db := getDB(DBFilepath)
	defer db.Close()

	rows, err := db.Query("SELECT source_host, created FROM monitor;")
	assert.Equal(t, nil, err)
	defer rows.Close()

	k := 0
	for rows.Next() {
		k++
		var sourceHost string
		var created string
		err = rows.Scan(&sourceHost, &created)
		assert.Equal(t, nil, err)
		assert.Equal(t, "http://a", sourceHost)
	}
	assert.Equal(t, 1, k)
}

func TestSaveRecordToMonitor_BadExternalLink(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)

	m := Monitor{}
	m.SourceHost = "http://a"
	m.ExternalLink = "http://somebaddomain.com/'.show_site_name($line['url'],100).'/"
	m.Count = 800
	m.ExternalHost = "b"
	res := SaveRecordToMonitor(DBFilepath, m)
	assert.Equal(t, true, res)

	db := getDB(DBFilepath)
	defer db.Close()

	rows, err := db.Query("SELECT external_link FROM monitor;")
	assert.Equal(t, nil, err)
	defer rows.Close()

	k := 0
	for rows.Next() {
		k++
		var fetchedExternalLink string
		err = rows.Scan(&fetchedExternalLink)
		assert.Equal(t, nil, err)
		assert.Equal(t, m.ExternalLink, fetchedExternalLink)
	}
	assert.Equal(t, 1, k)
}

func TestGetAllDataFromSqlite_MapToStruct(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)

	err := SaveHostType(DBFilepath, "host1", "type1")
	assert.NoError(t, err)

	err = SaveHostType(DBFilepath, "host2", "type2")
	assert.NoError(t, err)

	err = SaveHostType(DBFilepath, "host4", "type4")
	assert.NoError(t, err)

	for i := int(0); i < 10; i++ {
		m := Monitor{}
		m.SourceHost = "host1"
		m.ExternalLink = "http://b/1?" + strconv.Itoa(i)
		m.Count = 800 + i
		m.ExternalHost = "b"
		_ = SaveRecordToMonitor(DBFilepath, m)
	}

	m := Monitor{}
	m.SourceHost = "host2"
	m.ExternalLink = "http://b/1?10"
	m.Count = 810
	m.ExternalHost = "host1"
	_ = SaveRecordToMonitor(DBFilepath, m)

	m = Monitor{}
	m.SourceHost = "host3"
	m.ExternalLink = "http://b/1?10"
	m.Count = 810
	m.ExternalHost = "host3"
	_ = SaveRecordToMonitor(DBFilepath, m)

	m = Monitor{}
	m.SourceHost = "host3"
	m.ExternalLink = "http://b/1?10"
	m.Count = 810
	m.ExternalHost = "host4"
	_ = SaveRecordToMonitor(DBFilepath, m)

	monitors, err := GetAllDataFromMonitor(DBFilepath, 9)
	assert.NoError(t, err)
	assert.Equal(t, 13, len(monitors))
	assert.Equal(t, "http://b/1?0", monitors[0].ExternalLink)
	assert.Equal(t, "http://b/1?9", monitors[9].ExternalLink)

	assert.Equal(t, "type1", monitors[0].SourceHostType)
	assert.Equal(t, "N", monitors[0].ExternalHostType)
	assert.Equal(t, "type2", monitors[10].SourceHostType)
	assert.Equal(t, "type1", monitors[10].ExternalHostType)
	assert.Equal(t, "N", monitors[11].SourceHostType)
	assert.Equal(t, "N", monitors[11].ExternalHostType)
	assert.Equal(t, "type4", monitors[12].ExternalHostType)

	//for _, m := range monitors {
	//	log.Printf("[%v] [%v] %v %v %v", m.SourceHostType, m.ExternalHostType, m.Created, m.ExternalHost, m.SourceHost)
	//}
}

func TestParseSqliteDate(t *testing.T) {
	sqliteDate := "2016-12-13T07:17:23Z"
	parsedTime, _ := ParseSqliteDate(sqliteDate)
	assert.Equal(t, "2016", strconv.Itoa(parsedTime.Year()))
	assert.Equal(t, "12", strconv.Itoa(int(parsedTime.Month())))
	assert.Equal(t, "13", strconv.Itoa(parsedTime.Day()))
}

func TestGetDataFromMonitor__ByDays(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	err := SaveHostType(DBFilepath, "host1", "type1")
	assert.NoError(t, err)

	for i := int(0); i < 100; i++ {
		m := Monitor{}
		m.SourceHost = "http://a"
		m.ExternalLink = "http://b/1?" + strconv.Itoa(i)
		m.Count = 800 + i
		m.ExternalHost = "b"
		_ = SaveRecordToMonitor(DBFilepath, m)
	}

	dates, err := GetAllDaysFromMonitor(DBFilepath, "")
	assert.NoError(t, err)
	assert.Equal(t, len(dates), 1)

	monitors, err := GetAllDataFromMonitorByDay(DBFilepath, dates[0])
	assert.NoError(t, err)
	assert.Equal(t, len(monitors), 100)

	monitors, err = GetAllDataFromMonitorByDay(DBFilepath, "2122-01-01")
	assert.NoError(t, err)
	assert.Equal(t, len(monitors), 0)

	monitors, err = GetAllDataFromMonitorByDay(DBFilepath, "1812-01-01")
	assert.NoError(t, err)
	assert.Equal(t, len(monitors), 0)
}

func TestCrawlStatus(t *testing.T) {
	CreateDBIfNotExists(DBFilepath)

	SetCrawlStatus(DBFilepath, "Crawling...")
	s1, _ := GetCrawlStatus(DBFilepath)
	assert.Equal(t, "Crawling...", s1)

	SetCrawlStatus(DBFilepath, "Crawl done")
	GetCrawlStatus(DBFilepath)
	s2, _ := GetCrawlStatus(DBFilepath)
	assert.Equal(t, "Crawl done", s2)

	CreateDBIfNotExists(DBFilepath)
	SetCrawlStatus(DBFilepath, "Crawling...")
	s3, _ := GetCrawlStatus(DBFilepath)
	assert.Equal(t, "Crawling...", s3)
}

func TestSaveHostType(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	err := SaveHostType(DBFilepath, "host1", "type1")
	assert.NoError(t, err)
	err = SaveHostType(DBFilepath, "host1", "type1")
	assert.Error(t, err)
	err = SaveHostType(DBFilepath, "host1", "type2")
	assert.Error(t, err)
	err = SaveHostType(DBFilepath, "host2", "type1")
	assert.NoError(t, err)
}

func TestGetNewDataForDate(t *testing.T) {
	// fixture for different days
	// change SaveRecordToMonitor func to use date as a param, or create new method
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)

	monitorNoDate := Monitor{}
	monitorNoDate.SourceHost = "host2"
	monitorNoDate.ExternalLink = "http://b/1?10"
	monitorNoDate.Count = 810
	monitorNoDate.ExternalHost = "host0"
	_ = SaveRecordToMonitor(DBFilepath, monitorNoDate)

	monitorWithDate := Monitor{}
	monitorWithDate.SourceHost = "host2"
	monitorWithDate.ExternalLink = "http://b/1?10"
	monitorWithDate.Count = 810
	monitorWithDate.ExternalHost = "host1"
	monitorWithDate.Created = "2016-12-13"
	_ = SaveRecordToMonitor(DBFilepath, monitorWithDate)

	monitors, err := GetAllDataFromMonitor(DBFilepath, 0)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(monitors))
	assert.Equal(t, "host2", monitors[0].SourceHost)
	assert.Equal(t, "http://b/1?10", monitors[0].ExternalLink)
	assert.Equal(t, 810, monitors[0].Count)

	// lets parse the time to compare dates
	log.Print(monitors[0].Created)
	t1, _ := ParseMysqlDate(monitors[0].Created)
	t2 := time.Now().UTC()

	log.Printf("%v", t1)

	assert.True(t, (t1.Year() == t2.Year() && t1.YearDay() == t2.YearDay()))

	t1, _ = ParseSqliteDate(monitors[1].Created)
	t2 = time.Now()
	assert.False(t, (t1.Year() == t2.Year() && t1.YearDay() == t2.YearDay()))

	log.Print(t2.Format("2006-01-02"))

	monitors, err = GetNewExtractedHostsForDay(DBFilepath, t2.Format("2006-01-02"))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(monitors))
}

func TestGetNewDataByExternalHost(t *testing.T) {
	// fixture for different days
	// change SaveRecordToMonitor func to use date as a param, or create new method
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)

	monitorNoDate := Monitor{}
	monitorNoDate.SourceHost = "host2"
	monitorNoDate.ExternalLink = "http://b/1?10"
	monitorNoDate.Count = 810
	monitorNoDate.ExternalHost = "host0"
	_ = SaveRecordToMonitor(DBFilepath, monitorNoDate)

	monitorWithDate := Monitor{}
	monitorWithDate.SourceHost = "host2"
	monitorWithDate.ExternalLink = "http://b/1?10"
	monitorWithDate.Count = 811
	monitorWithDate.ExternalHost = "host1"
	monitorWithDate.Created = "2016-12-13"
	_ = SaveRecordToMonitor(DBFilepath, monitorWithDate)

	monitors, err := GetAllDataFromMonitorByExternalHost(DBFilepath, "host1")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(monitors))
	assert.Equal(t, 811, monitors[0].Count)
}

func TestDeleteTypesTable(t *testing.T) {
	err := DeleteTypesTable(DBFilepath)
	assert.NoError(t, err)
}

func TestGetAllTypes(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	err := SaveHostType(DBFilepath, "host1", "type1")
	assert.NoError(t, err)
	err = SaveHostType(DBFilepath, "host2", "type2")
	assert.NoError(t, err)
	hosts, err := GetAllTypes(DBFilepath)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(hosts))
}

func TestUpdateType(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)

	err := UpdateOrCreateHostType(DBFilepath, "test.com", "H")
	assert.NoError(t, err)
	types, err := GetAllTypes(DBFilepath)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(types))
	assert.Equal(t, "test.com", types[0].HostName)
	assert.Equal(t, "H", types[0].HostType)

	err = UpdateOrCreateHostType(DBFilepath, "test.com", "M")
	assert.NoError(t, err)
	types, err = GetAllTypes(DBFilepath)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(types))
	assert.Equal(t, "test.com", types[0].HostName)
	assert.Equal(t, "M", types[0].HostType)
}

func TestDeleteHost(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	err := SaveHostType(DBFilepath, "host1", "type1")
	assert.NoError(t, err)
	err = SaveHostType(DBFilepath, "host2", "type2")
	assert.NoError(t, err)
	err = DeleteHost(DBFilepath, "1")
	assert.NoError(t, err)
	hosts, err := GetAllTypes(DBFilepath)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(hosts))
}

func TestPerfomanceReport(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	for i := int(0); i < 5; i++ {
		m := Monitor{}
		m.Created = "2017-08-01"
		m.SourceHost = "host1"
		m.ExternalLink = "http://b/1?" + strconv.Itoa(i)
		m.Count = 20
		m.ExternalHost = "b"
		_ = SaveRecordToMonitor(DBFilepath, m)
	}
	for i := int(0); i < 5; i++ {
		m := Monitor{}
		m.Created = "2017-08-01"
		m.SourceHost = "host2"
		m.ExternalLink = "http://b/1?" + strconv.Itoa(i)
		m.Count = 20
		m.ExternalHost = "b"
		_ = SaveRecordToMonitor(DBFilepath, m)
	}
	for i := int(0); i < 5; i++ {
		m := Monitor{}
		m.Created = "2017-08-01"
		m.SourceHost = "host2"
		m.ExternalLink = "http://b/1?" + strconv.Itoa(i)
		m.Count = 20
		m.ExternalHost = "b"
		_ = SaveRecordToMonitor(DBFilepath, m)
	}
	for i := int(0); i < 5; i++ {
		m := Monitor{}
		m.Created = "2017-08-01"
		m.SourceHost = "host3"
		m.ExternalLink = "http://b/1?" + strconv.Itoa(i)
		m.Count = 20
		m.ExternalHost = "b"
		_ = SaveRecordToMonitor(DBFilepath, m)
	}
	for i := int(0); i < 5; i++ {
		m := Monitor{}
		m.Created = "2017-08-02"
		m.SourceHost = "host1"
		m.ExternalLink = "http://b/1?" + strconv.Itoa(i)
		m.Count = 20
		m.ExternalHost = "b"
		_ = SaveRecordToMonitor(DBFilepath, m)
	}
	for i := int(0); i < 5; i++ {
		m := Monitor{}
		m.Created = "2017-08-02"
		m.SourceHost = "host2"
		m.ExternalLink = "http://b/1?" + strconv.Itoa(i)
		m.Count = 20
		m.ExternalHost = "b"
		_ = SaveRecordToMonitor(DBFilepath, m)
	}
	for i := int(0); i < 5; i++ {
		m := Monitor{}
		m.Created = "2017-08-02"
		m.SourceHost = "host3"
		m.ExternalLink = "http://b/1?" + strconv.Itoa(i)
		m.Count = 20
		m.ExternalHost = "b"
		_ = SaveRecordToMonitor(DBFilepath, m)
	}
	for i := int(0); i < 5; i++ {
		m := Monitor{}
		m.Created = "2017-08-02"
		m.SourceHost = "host2"
		m.ExternalLink = "http://b/1?" + strconv.Itoa(i)
		m.Count = 20
		m.ExternalHost = "a"
		_ = SaveRecordToMonitor(DBFilepath, m)
	}
	for i := int(0); i < 5; i++ {
		m := Monitor{}
		m.Created = "2017-08-02"
		m.SourceHost = "host3"
		m.ExternalLink = "http://b/1?" + strconv.Itoa(i)
		m.Count = 10 + i
		m.ExternalHost = "b"
		_ = SaveRecordToMonitor(DBFilepath, m)
	}
	for i := int(0); i < 5; i++ {
		m := Monitor{}
		m.Created = "2017-08-03"
		m.SourceHost = "host2"
		m.ExternalLink = "http://b/1?" + strconv.Itoa(i)
		m.Count = 20
		m.ExternalHost = "a"
		_ = SaveRecordToMonitor(DBFilepath, m)
	}
	for i := int(0); i < 5; i++ {
		m := Monitor{}
		m.Created = "2017-08-03"
		m.SourceHost = "host2"
		m.ExternalLink = "http://b/1?" + strconv.Itoa(i)
		m.Count = 20
		m.ExternalHost = "b"
		_ = SaveRecordToMonitor(DBFilepath, m)
	}
	perfomanceData, err := PerfomanceReport(DBFilepath, "b")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(perfomanceData))
	assert.Equal(t, 400, perfomanceData[0].Count)
	assert.Equal(t, 3, perfomanceData[0].SourceHostCount)
}

func TestAddFeaturedHost(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	_, err := AddFeaturedHost(DBFilepath, "host1")
	assert.NoError(t, err)
	hosts, err := GetFeaturedHosts(DBFilepath)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(hosts))
}

func TestRemoveFeaturedHost(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	_, err := AddFeaturedHost(DBFilepath, "host1")
	assert.NoError(t, err)
	err = RemoveFeaturedHost(DBFilepath, "host1")
	assert.NoError(t, err)
	hosts, err := GetFeaturedHosts(DBFilepath)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(hosts))
}

func TestAddExistingFeaturedHost(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	_, err := AddFeaturedHost(DBFilepath, "host1")
	assert.NoError(t, err)
	_, err = AddFeaturedHost(DBFilepath, "host1")
	assert.NoError(t, err)
	hosts, err := GetFeaturedHosts(DBFilepath)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(hosts))
}

func TestGetFeaturedHosts(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	_, err := AddFeaturedHost(DBFilepath, "host1")
	assert.NoError(t, err)
	_, err = AddFeaturedHost(DBFilepath, "host2")
	assert.NoError(t, err)
	_, err = AddFeaturedHost(DBFilepath, "host3")
	assert.NoError(t, err)
	hosts, err := GetFeaturedHosts(DBFilepath)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(hosts))
}

func TestAddStopHost(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	err := AddStopHost(DBFilepath, "host1")
	assert.NoError(t, err)
	hosts, err := GetStopHosts(DBFilepath)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(hosts))
}

func TestRemoveStopHost(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	err := AddStopHost(DBFilepath, "host1")
	assert.NoError(t, err)
	err = RemoveStopHost(DBFilepath, "host1")
	assert.NoError(t, err)
	hosts, err := GetStopHosts(DBFilepath)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(hosts))
}

func TestGetStopHosts(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	err := AddStopHost(DBFilepath, "host1")
	assert.NoError(t, err)
	err = AddStopHost(DBFilepath, "host2")
	assert.NoError(t, err)
	err = AddStopHost(DBFilepath, "host3")
	assert.NoError(t, err)
	hosts, err := GetStopHosts(DBFilepath)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(hosts))
}

func TestSaveGrabberData(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	gd := GrabberData{}
	gd.Created = "2017-08-02"
	gd.Service = "Alexa"
	gd.Host = "namba.kg"
	gd.Data = "123"
	success := SaveGrabbedData(DBFilepath, gd)
	assert.Equal(t, true, success)
}

func TestPerfomanceReportByHostTypesNoErrors(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	_, err := PerfomanceReportByHostTypes(DBFilepath, "namba.kg")
	assert.NoError(t, err)
}

func TestPerfomanceReportGrabberDataNoErrors(t *testing.T) {
	TruncateDB(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	gd := GrabberData{}
	gd.Created = "2017-08-02"
	gd.Service = "Alexa"
	gd.Host = "namba.kg"
	gd.Data = "123"
	_ = SaveGrabbedData(DBFilepath, gd)
	data, err := PerfomanceReportGrabberData(DBFilepath, "alexa", "namba.kg")
	assert.NoError(t, err)
	fmt.Println(data)
	assert.Equal(t, 1, len(data))
	assert.Equal(t, "123", data[0].Data)
}
