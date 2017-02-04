package lib

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"fmt"
)

type Monitor struct {
	SourceHost string
	ExternalLink string
	Count int
	ExternalHost string
	Created string
	SourceHostType string
	ExternalHostType string
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

func GetAllDataFromMonitor(dbFilepath string, count int) ([]Monitor, error) {
	db, err := sql.Open("sqlite3", dbFilepath) // TODO: need to remove duplicates
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("SELECT m.source_host, m.external_link, m.count, m.external_host, m.created, " +
		"(CASE WHEN t.hostname != m.source_host THEN 'H' ELSE t.hosttype END) as 'source_host_type'," +
		"(CASE WHEN t.hostname != m.external_host THEN 'H' ELSE t.hosttype END) as 'external_host_type'" +
		"FROM monitor as m, types as t " +
		"WHERE m.count > %d;", count))
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

func GetAllDataFromMonitorByDay(dbFilepath string, day string) ([]Monitor, error) {
	db, err := sql.Open("sqlite3", dbFilepath) // TODO: need to remove duplicates
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT source_host, external_link, count, external_host, created FROM monitor WHERE created >= '%s' AND created <= date('%s', '+1 day');", day, day))
	if err != nil {
		log.Printf("Error getting data from monitor: %v", err)
		return nil, err
	}
	defer rows.Close()

	var data []Monitor
	for rows.Next() {
		m := Monitor{}
		err = rows.Scan(&m.SourceHost, &m.ExternalLink, &m.Count, &m.ExternalHost, &m.Created)
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

func DeleteTypesTable(dbFilepath string) error{
	db, err := sql.Open("sqlite3", dbFilepath) // TODO: need to remove duplicates
	if err != nil {
		log.Print(err)
		return err
	}
	defer db.Close()

	rows, err := db.Query("DELETE FROM types;")
	if err != nil {
		log.Print(err)
		return err
	}
	defer rows.Close()
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
