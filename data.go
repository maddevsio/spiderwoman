package main

import (
	"strings"
	"log"
	"github.com/maddevsio/spiderwoman/lib"
	"fmt"
)

func createXLS_BackupDB_Zip() {
	createAllXLSByDays()

	log.Print("Backuping database")
	err = lib.BackupDatabase(sqliteDBPath)
	if (err != nil) {
		log.Printf("Backup error: %v", err)
	} else {
		log.Print("Database has been copied to /tmp/res.db")
	}

	log.Print("Zip XLS File")
	err = lib.ZipFile(excelFilePath, excelZipFilePath)
	if (err != nil) {
		log.Printf("Zip error: %v", err)
	} else {
		log.Printf("Zipped xls file was saved in %v", excelZipFilePath)
	}
}

func createAllXLSByDays() {
	days, _ := lib.GetAllDaysFromMonitor(sqliteDBPath)
	for _, day := range days {
		log.Printf("Creating XLS file for day %v", day)
		lib.CreateExcelFromDB(sqliteDBPath, fmt.Sprintf(excelFilePath, day), day)
	}
}