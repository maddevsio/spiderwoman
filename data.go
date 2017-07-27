package main

import (
	"log"
	"github.com/maddevsio/spiderwoman/lib"
	"fmt"
)

func createAllXLSByDays(sqliteDBPath string) {
	days, _ := lib.GetAllDaysFromMonitor(sqliteDBPath)
	for _, day := range days {
		log.Printf("Creating XLS file for day %v", day)
		lib.CreateExcelFromDB(sqliteDBPath, fmt.Sprintf(excelFilePath, day), day)
	}
}