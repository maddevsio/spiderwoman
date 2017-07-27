package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"fmt"
)

func main() {
	dbSqlite, err := sql.Open("sqlite3", "../testdata/res.db")
	if err != nil {
		log.Panic(err)
	}
	defer dbSqlite.Close()

	dbMysql, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/spiderwoman?multiStatements=true")
	if err != nil {
		log.Panic(err)
	}
	err = dbMysql.Ping()
	if err != nil {
		log.Panic(err)
	}
	defer dbMysql.Close()

	///////

	query := fmt.Sprint("SELECT id FROM types")

	rows, err := dbSqlite.Query(query)

	if err != nil {
		log.Panicf("Error getting data from types: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		log.Print(id)
	}
}