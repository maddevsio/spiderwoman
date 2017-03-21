package main

import (
	"log"
	"github.com/maddevsio/spiderwoman/lib"
)

func initialize(path Path) {
	lib.CreateDBIfNotExists(path.SqliteDBPath)
	lib.ClearResolveCache()
	err = lib.PopulateHostsAndTypes(path.SqliteDBPath, path.SitesFilepath, path.SitesDefaultFilepath)
	if err != nil {
		log.Fatal("Types population error")
	}
}