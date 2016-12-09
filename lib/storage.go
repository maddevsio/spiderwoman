package lib

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func CreateDBIfNotExists(dbFilepath string) {
	db, err := sql.Open("sqlite3", dbFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table monitor (
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
