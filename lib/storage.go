package lib

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

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

	//query := fmt.Sprintf("insert into monitor(source_host, external_link, count, external_host, created) values('%v', '%v', %v, '%v', DateTime('now'))",
	//	source_host,
	//	external_link,
	//	count,
	//	external_host,
	//)

}

func GetAllDataFromMonitor(dbFilepath string) {

}

func ParseSqliteDate(sqliteDate string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05Z", sqliteDate)
}
