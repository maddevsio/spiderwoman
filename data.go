package main

import (
	"log"
	"github.com/maddevsio/spiderwoman/lib"
	"fmt"
)

func createXLS_BackupDB(sqliteDBPath string) {
	createAllXLSByDays(sqliteDBPath)

	log.Print("Backuping database")
	err = lib.BackupDatabase(sqliteDBPath)
	if (err != nil) {
		log.Printf("Backup error: %v", err)
	} else {
		log.Print("Database has been copied to /tmp/res.db")
	}
}

func createAllXLSByDays(sqliteDBPath string) {
	
	days, _ := lib.GetAllDaysFromMonitor(sqliteDBPath)
	for _, day := range days {
		log.Printf("Creating XLS file for day %v", day)
		lib.CreateExcelFromDB(sqliteDBPath, fmt.Sprintf(excelFilePath, day), day)
	}
}