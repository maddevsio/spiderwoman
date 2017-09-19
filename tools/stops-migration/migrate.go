package main

import (
	"github.com/maddevsio/spiderwoman/lib"
)

func main() {
	lib.CreateDBIfNotExists("spiderwoman")
	lib.MigrateStopHosts("spiderwoman", "../../stops.txt")
}