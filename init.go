package main

import (
	"github.com/maddevsio/spiderwoman/lib"
)

func initialize(path Path) {
	lib.CreateDBIfNotExistsAndMigrate(path.SqliteDBPath)
	lib.ClearResolveCache()
}
