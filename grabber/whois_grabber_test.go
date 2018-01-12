package grabber

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/maddevsio/spiderwoman/lib"
)

func TestWhoisValid(t *testing.T) {
	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExistsAndMigrate(dbName)
	whois := WhoisGrabber{Service{Name:"Whois"}}
	whoisData, err := whois.Do("google.com")
	assert.NoError(t, err)
	fmt.Print(whoisData)
}

func TestWhoisInvalid1(t *testing.T) {
	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExistsAndMigrate(dbName)
	whois := WhoisGrabber{Service{Name:"Whois"}}
	whoisData, err := whois.Do("non.existent.domain")
	assert.NoError(t, err)
	fmt.Print(whoisData)
}

func TestWhoisExotic1(t *testing.T) {
	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExistsAndMigrate(dbName)
	whois := WhoisGrabber{Service{Name:"Whois"}}
	whoisData, err := whois.Do("blah11blah22.travel")
	assert.NoError(t, err)
	fmt.Print(whoisData)
}

func TestWhoisExotic2(t *testing.T) {
	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExistsAndMigrate(dbName)
	whois := WhoisGrabber{Service{Name:"Whois"}}
	whoisData, err := whois.Do("blah11blah22.bid")
	assert.NoError(t, err)
	fmt.Print(whoisData)
}

func TestWhoisExotic3(t *testing.T) {
	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExistsAndMigrate(dbName)
	whois := WhoisGrabber{Service{Name:"Whois"}}
	whoisData, err := whois.Do("blah11blah22.website")
	assert.NoError(t, err)
	fmt.Print(whoisData)
}

func TestWhoisExotic4(t *testing.T) {
	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExistsAndMigrate(dbName)
	whois := WhoisGrabber{Service{Name:"Whois"}}
	whoisData, err := whois.Do("blah11blah22.online")
	assert.NoError(t, err)
	fmt.Print(whoisData)
}

func TestWhoisExotic5(t *testing.T) {
	lib.TruncateDB(dbName)
	lib.CreateDBIfNotExistsAndMigrate(dbName)
	whois := WhoisGrabber{Service{Name:"Whois"}}
	whoisData, err := whois.Do("blah11blah22.center")
	assert.NoError(t, err)
	fmt.Print(whoisData)
}