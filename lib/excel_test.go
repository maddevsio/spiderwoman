package lib

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestCreateExcel(t *testing.T) {
	os.Remove("./MyXLSXFile.xlsx")
	CreateDummySheet()
	_, err := os.Stat("./MyXLSXFile.xlsx");
	assert.Equal(t, nil, err)
}
