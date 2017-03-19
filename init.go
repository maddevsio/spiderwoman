package main

import (
	"log"
	"github.com/maddevsio/spiderwoman/lib"
)

func initialize() {
	lib.CreateDBIfNotExists(sqliteDBPath)
	lib.ClearResolveCache()
	err = lib.PopulateHostsAndTypes(sqliteDBPath, lib.SitesFilepath, lib.SitesDefaultFilepath)
	if err != nil {
		log.Fatal("Types population error")
	}
}