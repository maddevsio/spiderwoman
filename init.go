package main

import (
	"github.com/maddevsio/spiderwoman/lib"
)

func initialize(path Path) {
	lib.CreateDBIfNotExists(path.SqliteDBPath)
	lib.ClearResolveCache()
}
