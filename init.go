package main

import (
	"github.com/maddevsio/spiderwoman/lib"
)

func initialize(path Path) {
	lib.CreateDBIfNotExists(path.SqliteDBPath)
	lib.ClearResolveCache()
	//err = lib.PopulateHostsAndTypes(path.SqliteDBPath, path.TypesFilePath, path.TypesDefaultFilePath)
	//if err != nil {
	//	log.Fatal("Types population error")
	//}
}