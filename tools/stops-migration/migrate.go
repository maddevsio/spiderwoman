package main

import (
	"github.com/maddevsio/spiderwoman/lib"
)

func main() {
	lib.CreateDBIfNotExistsAndMigrate("spiderwoman")
	lib.MigrateStopHosts("spiderwoman", "../../stops.txt")
}
