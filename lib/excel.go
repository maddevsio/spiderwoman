package lib


import (
	"github.com/tealeg/xlsx"
	"strconv"
	"log"
)

func CreateExcelFromDB(dbFilepath string, excelFilePath string) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Full Data")
	if err != nil {
		log.Print(err)
	}

	monitors, _ := GetAllDataFromMonitor(dbFilepath, 0)
	fillTheSheet(sheet, monitors)

	err = file.Save(excelFilePath)
	if err != nil {
		log.Print(err)
	}
}

func CreateEmptyExcel(excelFilePath string) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Empty")
	if err != nil {
		log.Print(err)
	}

	row := sheet.AddRow()
	cell := row.AddCell()
	cell.Value = ""

	err = file.Save(excelFilePath)
	if err != nil {
		log.Print(err)
	}
}

func AppendExcelFromDB(dbFilepath string, excelFilePath string, date string) error {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var err error

	file, err = xlsx.OpenFile(excelFilePath)
	if err != nil {
		log.Print(err)
		return err

	}
	sheet, err = file.AddSheet(date)
	if err != nil {
		log.Print(err)
		return err
	}

	monitors, _ := GetAllDataFromMonitorByDay(dbFilepath, date)
	fillTheSheet(sheet, monitors)

	err = file.Save(excelFilePath)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func fillTheSheet(sheet *xlsx.Sheet, monitors []Monitor) {
	row := sheet.AddRow()
	cell1 := row.AddCell()
 	cell1.Value = "Date"

	cell2 := row.AddCell()
	cell2.Value = "Source Host"

	cell3 := row.AddCell()
	cell3.Value = "Type SH"

	cell4 := row.AddCell()
	cell4.Value = "External Host"

	cell5 := row.AddCell()
	cell5.Value = "Type EH"

	cell6 := row.AddCell()
	cell6.Value = "Count"

	cell7 := row.AddCell()
	cell7.Value = "External Link"

	for _, monitor := range monitors {
		row := sheet.AddRow()

		cell1 := row.AddCell()
		cell1.Value = monitor.Created

		cell2 := row.AddCell()
		cell2.Value = monitor.SourceHost

		cell3 := row.AddCell()
		cell3.Value = monitor.SourceHostType

		cell4 := row.AddCell()
		cell4.Value = monitor.ExternalHost

		cell5 := row.AddCell()
		cell5.Value = monitor.ExternalHostType

		cell6 := row.AddCell()
		cell6.Value = strconv.Itoa(monitor.Count)

		cell7 := row.AddCell()
		cell7.Value = monitor.ExternalLink
	}
}