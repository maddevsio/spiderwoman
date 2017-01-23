package lib

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
)

func TestCreateExcel(t *testing.T) {
	dbFilePath := "/tmp/spiderwoman.db"
	excelFilePath := "/tmp/spiderwoman.xls"

	os.Remove(excelFilePath)
	CreateExcelFromDB(dbFilePath, excelFilePath)
	_, err := os.Stat(excelFilePath);

	assert.Equal(t, nil, err)
}

func TestAppendExcel(t *testing.T) {
	dbFilePath := "/tmp/spiderwoman.db"
	excelFilePath := "/tmp/spiderwoman.xls"

	AppendExcelFromDB(dbFilePath, excelFilePath, "2017-01-20")
	_, err := os.Stat(excelFilePath);

	assert.Equal(t, nil, err)
}

func TestAppendExcelNoFile(t *testing.T) {
	dbFilePath := "/tmp/spiderwoman.db"
	excelFilePath := "/tmp/spiderwoman.xls"

	os.Remove(excelFilePath)
	err := AppendExcelFromDB(dbFilePath, excelFilePath, "2017-01-20")

	assert.Error(t, err)

	CreateEmptyExcel(excelFilePath)
	err = AppendExcelFromDB(dbFilePath, excelFilePath, "2017-01-20")
	assert.NoError(t, err)

	_, err = os.Stat(excelFilePath);
	assert.NoError(t, err)
}



