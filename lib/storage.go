package lib

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Monitor struct {
	ID               int
	SourceHost       string
	ExternalLink     string
	Count            int
	ExternalHost     string
	Created          string
	SourceHostType   string
	ExternalHostType string
}

type HostItem struct {
	ID       int64
	HostName string
	HostType string
}

type Hosts struct {
	Host []HostItem
}

type PerfomanceReportResponse struct {
	Created         string
	SourceHostCount int
	Count           int
}

func TruncateDB(dbName string) {
	db := getDB(dbName)
	defer db.Close()
	sqlStmt := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`; CREATE DATABASE `%s`;", dbName, dbName)
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func getDB(dbName string) *sql.DB {
	db, err := sql.Open("mysql", "root@tcp(mysql:3306)/"+dbName+"?multiStatements=true")
	if err != nil {
		log.Printf("===%v===", err)
		log.Panic(err)
	}
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	return db
}

func CreateDBIfNotExists(dbFilepath string) {
	db := getDB(dbFilepath)
	defer db.Close()

	sqlStmt := `
	create table if not exists monitor (
		ID int NOT NULL AUTO_INCREMENT,
        PRIMARY KEY (id),
		source_host varchar(255),
		external_link varchar(255),
		count int,
		external_host varchar(255),
		created date
	) DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci ENGINE=InnoDB;
    		create table if not exists status (
		ID int NOT NULL AUTO_INCREMENT,
	        PRIMARY KEY (id),
		status_key varchar(255),
		status_value varchar(255)
	);
    	insert into status(status_key, status_value) values('crawl', 'Crawl done');
    	create table if not exists types (
		ID int NOT NULL AUTO_INCREMENT,
        	PRIMARY KEY (id),
		hostname varchar(255),
		hosttype varchar(255),
		CONSTRAINT hostname_uniq UNIQUE (hostname)
	) DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci ENGINE=InnoDB;
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func SaveRecordToMonitor(dbFilepath string, monitor Monitor) bool {
	db := getDB(dbFilepath)
	defer db.Close()

	stmt, err := db.Prepare("insert into monitor(source_host, external_link, count, external_host, created) values(?, ?, ?, ?, NOW())")
	if monitor.Created != "" {
		stmt, err = db.Prepare("insert into monitor(source_host, external_link, count, external_host, created) values(?, ?, ?, ?, ?)")
	}
	if err != nil {
		log.Print(err)
		return false
	}

	_, err = stmt.Exec(monitor.SourceHost, monitor.ExternalLink, monitor.Count, monitor.ExternalHost)
	if monitor.Created != "" {
		_, err = stmt.Exec(monitor.SourceHost, monitor.ExternalLink, monitor.Count, monitor.ExternalHost, monitor.Created)
	}
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}

func GetAllDataFromMonitor(dbFilepath string, count int) ([]Monitor, error) {
	db := getDB(dbFilepath)
	defer db.Close()

	sql := fmt.Sprintf("SELECT m.id, m.source_host, m.external_link, m.count, m.external_host, m.created, "+
		"coalesce(t1.hosttype,'N') as 'source_host_type', "+
		"coalesce(t2.hosttype,'N') as 'external_host_type' "+
		"FROM monitor as m "+
		"LEFT OUTER JOIN types as t1 ON t1.hostname=m.source_host "+
		"LEFT OUTER JOIN types as t2 ON t2.hostname=m.external_host "+
		"WHERE m.count > %d ORDER BY id ASC", count)

	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("Error getting data from monitor: %v", err)
		return nil, err
	}
	defer rows.Close()

	var data []Monitor
	for rows.Next() {
		m := Monitor{}
		err = rows.Scan(&m.ID, &m.SourceHost, &m.ExternalLink, &m.Count, &m.ExternalHost, &m.Created, &m.SourceHostType, &m.ExternalHostType)
		data = append(data, m)
	}

	return data, nil
}

func UpdateOrCreateHostType(dbFilepath string, hostName string, hostType string) error {
	db := getDB(dbFilepath)
	defer db.Close()

	stmt, err := db.Prepare("REPLACE INTO types VALUES (NULL, ?, ?);")
	if err != nil {
		log.Print(err)
		return err
	}
	_, err = stmt.Exec(hostName, hostType)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func GetAllDataFromMonitorByExternalHost(dbFilepath string, host string) ([]Monitor, error) {
	db := getDB(dbFilepath)
	defer db.Close()

	query := fmt.Sprintf("SELECT m.id, m.source_host, m.external_link, m.count, m.external_host, m.created, "+
		"coalesce(t1.hosttype,'N') as 'source_host_type', "+
		"coalesce(t2.hosttype,'N') as 'external_host_type' "+
		"FROM monitor as m "+
		"LEFT OUTER JOIN types as t1 ON t1.hostname=m.source_host "+
		"LEFT OUTER JOIN types as t2 ON t2.hostname=m.external_host "+
		"WHERE m.external_host = '%s';", host)

	rows, err := db.Query(query)

	if err != nil {
		log.Printf("Error getting data from monitor: %v", err)
		return nil, err
	}
	defer rows.Close()

	var data []Monitor
	for rows.Next() {
		m := Monitor{}
		err = rows.Scan(&m.ID, &m.SourceHost, &m.ExternalLink, &m.Count, &m.ExternalHost, &m.Created, &m.SourceHostType, &m.ExternalHostType)
		data = append(data, m)
	}

	return data, nil
}

func GetAllDataFromMonitorByDay(dbFilepath string, day string) ([]Monitor, error) {
	db := getDB(dbFilepath)
	defer db.Close()

	query := fmt.Sprintf("SELECT m.id, m.source_host, m.external_link, m.count, m.external_host, m.created, "+
		"coalesce(t1.hosttype,'N') as 'source_host_type', "+
		"coalesce(t2.hosttype,'N') as 'external_host_type' "+
		"FROM monitor as m "+
		"LEFT OUTER JOIN types as t1 ON t1.hostname=m.source_host "+
		"LEFT OUTER JOIN types as t2 ON t2.hostname=m.external_host "+
		"WHERE m.created >= '%s' AND m.created <= ('%s' + INTERVAL 1 DAY);", day, day)

	rows, err := db.Query(query)

	if err != nil {
		log.Printf("Error getting data from monitor: %v", err)
		return nil, err
	}
	defer rows.Close()

	var data []Monitor
	for rows.Next() {
		m := Monitor{}
		err = rows.Scan(&m.ID, &m.SourceHost, &m.ExternalLink, &m.Count, &m.ExternalHost, &m.Created, &m.SourceHostType, &m.ExternalHostType)
		data = append(data, m)
	}

	return data, nil
}

func GetNewExtractedHostsForDay(dbFilepath string, day string) ([]Monitor, error) {
	db := getDB(dbFilepath)
	defer db.Close()

	query := fmt.Sprintf("SELECT m.id, m.source_host, m.external_link, m.count, m.external_host, m.created, "+
		"coalesce(t1.hosttype,'N') as 'source_host_type', "+
		"coalesce(t2.hosttype,'N') as 'external_host_type' "+
		"FROM monitor as m "+
		"LEFT OUTER JOIN types as t1 ON t1.hostname=m.source_host "+
		"LEFT OUTER JOIN types as t2 ON t2.hostname=m.external_host "+
		"WHERE m.external_host not in (select distinct external_host from monitor where created < '%s') and "+
		"m.created >= '%s' AND m.created <= ('%s' + INTERVAL 1 DAY);", day, day, day)

	rows, err := db.Query(query)

	if err != nil {
		log.Printf("Error getting data from monitor: %v", err)
		return nil, err
	}
	defer rows.Close()

	var data []Monitor
	for rows.Next() {
		m := Monitor{}
		err = rows.Scan(&m.ID, &m.SourceHost, &m.ExternalLink, &m.Count, &m.ExternalHost, &m.Created, &m.SourceHostType, &m.ExternalHostType)
		data = append(data, m)
	}

	return data, nil
}

func SetCrawlStatus(dbFilepath string, status string) bool {
	db := getDB(dbFilepath)
	defer db.Close()

	stmt, err := db.Prepare("UPDATE status SET status_value=? WHERE status_key=?")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(status, "crawl")
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func GetCrawlStatus(dbFilepath string) (string, error) {
	db := getDB(dbFilepath)
	defer db.Close()

	rows, err := db.Query("SELECT status_value FROM status WHERE status_key='crawl';")
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	defer rows.Close()

	var status string

	for rows.Next() {
		err = rows.Scan(&status)
	}
	return status, nil
}

func GetAllDaysFromMonitor(dbFilepath string) ([]string, error) {
	db := getDB(dbFilepath)
	defer db.Close()

	rows, err := db.Query("SELECT DISTINCT created as mon FROM monitor ORDER BY created DESC;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var dates []string

	for rows.Next() {
		var date string
		err = rows.Scan(&date)
		if err == nil {
			dates = append(dates, date)
		} else {
			log.Printf("Error getting dates: %v", err)
		}
	}
	return dates, nil
}

func DeleteTypesTable(dbFilepath string) error {
	db := getDB(dbFilepath)
	defer db.Close()

	_, err := db.Exec("DELETE FROM types;")
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func SaveHostType(dbFilepath string, hostName string, hostType string) error {
	db := getDB(dbFilepath)
	defer db.Close()

	stmt, err := db.Prepare("insert into types(hostname, hosttype) values(?, ?)")
	if err != nil {
		log.Print(err)
		return err
	}

	_, err = stmt.Exec(hostName, hostType)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func ParseSqliteDate(sqliteDate string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05Z", sqliteDate)
}

func ParseMysqlDate(sqliteDate string) (time.Time, error) {
	return time.Parse("2006-01-02", sqliteDate)
}

func GetAllTypes(dbFilepath string) ([]HostItem, error) {
	db := getDB(dbFilepath)
	defer db.Close()

	query := fmt.Sprintf("SELECT t.id, t.hostname, t.hosttype FROM types as t")

	rows, err := db.Query(query)

	if err != nil {
		log.Printf("Error getting data from types: %v", err)
		return nil, err
	}
	defer rows.Close()

	var data []HostItem
	for rows.Next() {
		t := HostItem{}
		err = rows.Scan(&t.ID, &t.HostName, &t.HostType)
		data = append(data, t)
	}

	return data, nil
}

func GetUniqueTypes(dbFilepath string) ([]string, error) {
	db := getDB(dbFilepath)
	defer db.Close()

	query := fmt.Sprintf("SELECT DISTINCT t.hosttype FROM types as t")

	rows, err := db.Query(query)

	if err != nil {
		log.Printf("Error getting data from types: %v", err)
		return nil, err
	}
	defer rows.Close()

	var types []string
	for rows.Next() {
		var htype string
		err = rows.Scan(&htype)
		if err == nil {
			types = append(types, htype)
		} else {
			log.Printf("Error getting types: %v", err)
		}
	}
	return types, nil
}

func DeleteHost(dbFilepath string, hostID string) error {
	db := getDB(dbFilepath)
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM types WHERE id = ?")
	if err != nil {
		log.Print(err)
		return err
	}
	_, err = stmt.Exec(hostID)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func PerfomanceReport(dbFilepath string, host string) ([]PerfomanceReportResponse, error) {
	db := getDB(dbFilepath)
	defer db.Close()

	query := fmt.Sprintf("SELECT created, SUM(count), COUNT(DISTINCT `source_host`) as source_host_count "+
		"FROM monitor "+
		"WHERE external_host = '%s' GROUP BY created;", host)
	fmt.Println(query)

	rows, err := db.Query(query)

	if err != nil {
		log.Printf("Error getting data from monitor: %v", err)
		return nil, err
	}
	defer rows.Close()

	var data []PerfomanceReportResponse
	for rows.Next() {
		m := PerfomanceReportResponse{}
		err = rows.Scan(&m.Created, &m.SourceHostCount, &m.Count)
		data = append(data, m)
	}

	return data, nil
}
