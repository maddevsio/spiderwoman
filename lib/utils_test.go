package lib

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestMain(m *testing.M) {
	ClearResolveCache()
	os.Exit(m.Run())
}

func TestResolveOneRedirect(t *testing.T) {
	res := Resolve("http://bit.ly/2hcXx5Z", "http://bit.ly/", 10, false, "Googlebot", nil)
	assert.Equal(t, "https://maddevs.io/", res)
}

func TestResolveTwoRedirects(t *testing.T) {
	res := Resolve("http://ow.ly/Pkxu306YmRs", "http://ow.ly/", 10, false, "Googlebot", nil)
	assert.Equal(t, "https://maddevs.io/", res)
}

func TestResolveSSL(t *testing.T) {
	res := Resolve("https://bit.ly/2hcXx5Z", "http://bit.ly/", 10, false, "Googlebot", nil)
	assert.Equal(t, "https://maddevs.io/", res)

	res = Resolve("https://bit.ly/ID7AM5", "http://bit.ly/", 10, false, "Googlebot", nil)
	assert.Equal(t, "https://www.youtube.com/", res)
}

func TestResolveCache(t *testing.T) {
	res := Resolve("http://bit.ly/ItaROu", "http://bit.ly/", 10, false, "Googlebot", nil)
	assert.Equal(t, "https://duckduckgo.com/", res)
	assert.Equal(t, lastCachedReturn, false)
	assert.Equal(t, resolveCache["http://bit.ly/ItaROu"], "https://duckduckgo.com/")

	res = Resolve("http://bit.ly/ItaROu", "http://bit.ly/", 10, false, "Googlebot", nil)
	assert.Equal(t, "https://duckduckgo.com/", res)
	assert.Equal(t, lastCachedReturn, true)
}

func TestBackup(t *testing.T) {
	os.Remove("/tmp/res.db")
	err := BackupDatabase(DBFilepath)
	assert.NoError(t, err)
	_, err = os.Stat("/tmp/res.db")
	assert.Equal(t, false, os.IsNotExist(err))
}

func TestGetHostsFromFile(t *testing.T) {
	hosts, err := GetHostsFromFile("", "../sites.default.txt")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(hosts))
	assert.Equal(t, "nambataxi.kg", hosts[0])
	assert.Equal(t, "nambafood.kg", hosts[1])
}

func TestPopulateHostsAndTypes(t *testing.T) {
	os.Remove(DBFilepath)
	CreateDBIfNotExists(DBFilepath)
	err := PopulateHostsAndTypes(DBFilepath, "", "../sites.default.txt")
	assert.NoError(t, err)
	err = PopulateHostsAndTypes(DBFilepath, "", "../sites.default.txt")
	assert.NoError(t, err)
}

func TestMkdirAll(t *testing.T) {
	err := os.MkdirAll("/tmp/xls", 0777)
	assert.NoError(t, err)
}
