package grabber

import (
	"fmt"
	"testing"

	"github.com/maddevsio/spiderwoman/lib"
	"github.com/negah/alexa"
	"github.com/stretchr/testify/assert"
)

func TestAlexaLib(t *testing.T) {
	hosts := []string{"apteka312.kg", "apteka", "12341234"}
	for _, h := range hosts {
		rank, err := alexa.GlobalRank(h)
		// TODO: Not sure how to handle error from alexa lib properly.
		if err != nil {
			rank = "No rank"
			fmt.Println(err)
		}
		assert.Equal(t, "No rank", rank)
	}
}

func TestGrabberSavedNoRankAlexaData(t *testing.T) {
	alexaGrabber := AlexaGrabber{Service{Name: "Alexa"}}
	hosts := []string{"apteka312.kg", "apteka", "12341234"}
	for _, h := range hosts {
		success, err := GrabAndSave(alexaGrabber, h, dbName)
		if err != nil {
			fmt.Println(err)
		}
		d, err := lib.PerfomanceReportLatestGrabberData(dbName, "alexa", h)
		if err != nil {
			fmt.Println(err)
		}
		assert.NoError(t, err)
		assert.Equal(t, true, success)
		assert.Equal(t, "No rank", d.Data)
	}
}
