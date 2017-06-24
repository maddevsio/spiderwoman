package lib

import (
	"database/sql"
	"log"
	"time"

	"fmt"

	_ "github.com/mattn/go-sqlite3"
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

func CreateDBIfNotExists(dbFilepath string) {
	db, err := sql.Open("sqlite3", dbFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table if not exists monitor (
		id integer not null primary key,
		source_host text,
		external_link text,
		count int,
		external_host text,
		created date
	);
	create table if not exists status (
		id integer not null primary key,
		status_key text,
		status_value text
	);
	insert into status(status_key, status_value) values('crawl', 'Crawl done');
	create table if not exists types (
		id integer not null primary key,
		hostname text,
		hosttype text,
		CONSTRAINT hostname_uniq UNIQUE (hostname)
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func SaveRecordToMonitor(dbFilepath string, source_host string, external_link string, count int, external_host string) bool {
	db, err := sql.Open("sqlite3", dbFilepath)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into monitor(source_host, external_link, count, external_host, created) values(?, ?, ?, ?, DateTime('now'))")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(source_host, external_link, count, external_host)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func SaveRecordToMonitorStruct(dbFilepath string, monitor Monitor) bool {
	db, err := sql.Open("sqlite3", dbFilepath)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into monitor(source_host, external_link, count, external_host, created) values(?, ?, ?, ?, DateTime('now'))")
	if monitor.Created != "" {
		stmt, err = db.Prepare("insert into monitor(source_host, external_link, count, external_host, created) values(?, ?, ?, ?, ?)")
	}
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(monitor.SourceHost, monitor.ExternalLink, monitor.Count, monitor.ExternalHost)
	if monitor.Created != "" {
		_, err = stmt.Exec(monitor.SourceHost, monitor.ExternalLink, monitor.Count, monitor.ExternalHost, monitor.Created)
	}
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func GetAllDataFromMonitor(dbFilepath string, count int) ([]Monitor, error) {
	db, err := sql.Open("sqlite3", dbFilepath) // TODO: need to remove duplicates
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("SELECT m.id, m.source_host, m.external_link, m.count, m.external_host, m.created, "+
		"coalesce(t1.hosttype,'N') as 'source_host_type', "+
		"coalesce(t2.hosttype,'N') as 'external_host_type' "+
		"FROM monitor as m "+
		"LEFT OUTER JOIN types as t1 ON t1.hostname=m.source_host "+
		"LEFT OUTER JOIN types as t2 ON t2.hostname=m.external_host "+
		"WHERE m.count > %d;", count))
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
	db, err := sql.Open("sqlite3", dbFilepath) // TODO: need to remove duplicates
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT OR REPLACE INTO types VALUES (NULL, ?, ?);")
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
	db, err := sql.Open("sqlite3", dbFilepath) // TODO: need to remove duplicates
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
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
	db, err := sql.Open("sqlite3", dbFilepath) // TODO: need to remove duplicates
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT m.id, m.source_host, m.external_link, m.count, m.external_host, m.created, "+
		"coalesce(t1.hosttype,'N') as 'source_host_type', "+
		"coalesce(t2.hosttype,'N') as 'external_host_type' "+
		"FROM monitor as m "+
		"LEFT OUTER JOIN types as t1 ON t1.hostname=m.source_host "+
		"LEFT OUTER JOIN types as t2 ON t2.hostname=m.external_host "+
		"WHERE m.created >= '%s' AND m.created <= date('%s', '+1 day');", day, day)

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
	db, err := sql.Open("sqlite3", dbFilepath) // TODO: need to remove duplicates
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT m.source_host, m.external_link, m.count, m.external_host, m.created, "+
		"coalesce(t1.hosttype,'N') as 'source_host_type', "+
		"coalesce(t2.hosttype,'N') as 'external_host_type' "+
		"FROM monitor as m "+
		"LEFT OUTER JOIN types as t1 ON t1.hostname=m.source_host "+
		"LEFT OUTER JOIN types as t2 ON t2.hostname=m.external_host "+
		"WHERE m.external_host not in (select distinct external_host from monitor where created < '%s') and "+
		"m.created >= '%s' AND m.created <= date('%s', '+1 day');", day, day, day)

	rows, err := db.Query(query)

	if err != nil {
		log.Printf("Error getting data from monitor: %v", err)
		return nil, err
	}
	defer rows.Close()

	var data []Monitor
	for rows.Next() {
		m := Monitor{}
		err = rows.Scan(&m.SourceHost, &m.ExternalLink, &m.Count, &m.ExternalHost, &m.Created, &m.SourceHostType, &m.ExternalHostType)
		data = append(data, m)
	}

	return data, nil
}

func SetCrawlStatus(dbFilepath string, status string) bool {
	db, err := sql.Open("sqlite3", dbFilepath)
	if err != nil {
		log.Fatal(err)
		return false
	}
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
	db, err := sql.Open("sqlite3", dbFilepath) // TODO: need to remove duplicates
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer db.Close()

	rows, err := db.Query("SELECT status_value FROM status WHERE status_key='crawl';")
	defer rows.Close()

	var status string

	for rows.Next() {
		err = rows.Scan(&status)
	}
	return status, nil
}

func GetAllDaysFromMonitor(dbFilepath string) ([]string, error) {
	db, err := sql.Open("sqlite3", dbFilepath) // TODO: need to remove duplicates
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT DISTINCT strftime('%Y-%m-%d', created) as mon FROM monitor ORDER BY created DESC;")
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
	db, err := sql.Open("sqlite3", dbFilepath) // TODO: need to remove duplicates
	if err != nil {
		log.Print(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM types;")
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func SaveHostType(dbFilepath string, hostName string, hostType string) error {
	db, err := sql.Open("sqlite3", dbFilepath)
	if err != nil {
		log.Print(err)
		return err
	}
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

func GetAllTypes(dbFilepath string) ([]HostItem, error) {
	db, err := sql.Open("sqlite3", dbFilepath) // TODO: need to remove duplicates
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
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
	db, err := sql.Open("sqlite3", dbFilepath) // TODO: need to remove duplicates
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
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
	db, err := sql.Open("sqlite3", dbFilepath) // TODO: need to remove duplicates
	if err != nil {
		log.Fatal(err)
		return err
	}
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
