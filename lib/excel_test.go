package lib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateExcel(t *testing.T) {
	excelFilePath := "/tmp/spiderwoman.xls"

	TruncateDB(DBFilepath)
	CreateDBIfNotExistsAndMigrate(DBFilepath)
	CreateExcelFromDB(DBFilepath, excelFilePath, "2017-01-20")
	_, err := os.Stat(excelFilePath)

	assert.Equal(t, nil, err)
}

func TestCreateExcelNoFile(t *testing.T) {

}

func TestCreateExcelForOneDay(t *testing.T) {

}
