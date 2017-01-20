package lib


import (
	"fmt"
	"github.com/tealeg/xlsx"
	"strconv"
)

func CreateExcelFromDB(dbFilepath string, excelFilePath string) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}

	monitors, _ := GetAllDataFromMonitor(dbFilepath, 0)
	for _, monitor := range monitors {
		row := sheet.AddRow()

		cell1 := row.AddCell()
		cell1.Value = monitor.SourceHost

		cell2 := row.AddCell()
		cell2.Value = monitor.ExternalHost

		cell3 := row.AddCell()
		cell3.Value = monitor.ExternalLink

		cell4 := row.AddCell()
		cell4.Value = strconv.Itoa(monitor.Count)

		cell5 := row.AddCell()
		cell5.Value = monitor.Created
	}

	err = file.Save(excelFilePath)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func AppendExcelFromDB(dbFilepath string, excelFilePath string, date string) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var err error

	file, err = xlsx.OpenFile(excelFilePath)
	if err != nil {
		fmt.Printf(err.Error())
	}
	sheet, err = file.AddSheet(date)
	if err != nil {
		fmt.Printf(err.Error())
	}

	monitors, _ := GetAllDataFromMonitor(dbFilepath, 0)
	for _, monitor := range monitors {
		row := sheet.AddRow()

		cell1 := row.AddCell()
		cell1.Value = monitor.SourceHost

		cell2 := row.AddCell()
		cell2.Value = monitor.ExternalHost

		cell3 := row.AddCell()
		cell3.Value = monitor.ExternalLink

		cell4 := row.AddCell()
		cell4.Value = strconv.Itoa(monitor.Count)

		cell5 := row.AddCell()
		cell5.Value = monitor.Created
	}

	err = file.Save(excelFilePath)
	if err != nil {
		fmt.Printf(err.Error())
	}
}