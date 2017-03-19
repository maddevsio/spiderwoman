package main

import (
	"strings"
	"log"
	"github.com/maddevsio/spiderwoman/lib"
)

func createXLS_BackupDB_Zip() {
	days, _ := lib.GetAllDaysFromMonitor(sqliteDBPath)
	log.Printf("Appendig XLS file with sheet %v", days[0])
	err = lib.AppendExcelFromDB(sqliteDBPath, excelFilePath, days[0])
	if (err != nil && strings.Contains(err.Error(), "no such file or directory")) {
		lib.CreateEmptyExcel(excelFilePath)
		log.Print("Trying to create all sheets in excel file")
		for _, day := range days {
			log.Printf("Appendig XLS file with sheet %v", day)
			err = lib.AppendExcelFromDB(sqliteDBPath, excelFilePath, day)
			if err != nil {
				log.Print(err)
			}
		}
	}

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