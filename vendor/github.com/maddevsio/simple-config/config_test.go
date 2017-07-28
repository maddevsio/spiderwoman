package simple_config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigGet_Exist(t *testing.T) {
	config := NewSimpleConfig("./config.test", "yml")
	value := config.Get("testkey")
	assert.Equal(t, value, "test value")
}

func TestConfigGet_NonExist(t *testing.T) {
	config := NewSimpleConfig("./config.test", "yml")
	value := config.Get("allauakbarrr")
	assert.Nil(t, value)
}

func TestConfigGet_String(t *testing.T) {
	config := NewSimpleConfig("./config.test", "yml")
	value := config.GetString("testkey")
	assert.Equal(t, value, "test value")
}

func TestConfigGet_Env(t *testing.T) {
	os.Setenv("FROMENV", "from env")
	config := NewSimpleConfig("./config.test", "yml")
	value := config.GetString("fromenv")
	assert.Equal(t, value, "from env")
}
