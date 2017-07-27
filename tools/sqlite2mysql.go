package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"fmt"
	"strings"
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

	/////// types

	//query := fmt.Sprint("SELECT hostname, hosttype FROM types")
	//rows, err := dbSqlite.Query(query)
	//if err != nil {
	//	log.Panicf("Error getting data from types: %v", err)
	//}
	//defer rows.Close()
	//
	//for rows.Next() {
	//	var hostname string
	//	var hosttype string
	//	err = rows.Scan(&hostname, &hosttype)
	//	log.Printf("type: %s, host: %s", hosttype, hostname)
	//	stmt, err := dbMysql.Prepare("insert into types(hostname, hosttype) values(?, ?)")
	//	if err != nil {
	//		log.Print(err)
	//	}
	//	_, err = stmt.Exec(hostname, hosttype)
	//	if err != nil {
	//		log.Print(err)
	//	}
	//}

	/////// monitors

	query := fmt.Sprint("SELECT id, source_host, external_link, count, external_host, created FROM monitor ORDER BY id DESC")
	rows, err := dbSqlite.Query(query)
	if err != nil {
		log.Panicf("Error getting data from mons: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var source_host string
		var external_link string
		var count int
		var external_host string
		var created string
		err = rows.Scan(&id, &source_host, &external_link, &count, &external_host, &created)
		log.Printf("id %d", id)
		stmt, err := dbMysql.Prepare("insert into monitor(source_host, external_link, count, external_host, created) values(?, ?, ?, ?, ?)")
		if err != nil {
			log.Print(err)
		}

		created = strings.Split(created, "T")[0]

		_, err = stmt.Exec(source_host, external_link, count, external_host, created)
		if err != nil {
			log.Print(err)
		}
		stmt.Close()
	}
}