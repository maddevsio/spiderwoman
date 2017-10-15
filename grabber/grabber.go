package grabber

import (
	"fmt"
	"github.com/maddevsio/spiderwoman/lib"
)

type Grabber interface {
	CheckConnection() (int, error)
	GetServiceInfo() Service
	Do(featuredHost string) (string, error)
}

func GrabAndSave(g Grabber, featuredHost string, dbName string) (bool, error) {
	conn, err := g.CheckConnection()
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	if conn == 200 {
		rawData, err := g.Do(featuredHost)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		fmt.Println(rawData)
		gd := lib.GrabberData{}

		gd.Service = g.GetServiceInfo().Name
		gd.Host = featuredHost
		gd.Data = rawData
		success := lib.SaveGrabbedData(dbName, gd)
		return success, nil
	}
	return false, nil
}