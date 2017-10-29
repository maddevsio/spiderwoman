package grabber

import (
	"fmt"

	"github.com/negah/alexa"
)

type AlexaGrabber struct {
	Service
}

func (ag AlexaGrabber) CheckConnection() (int, error) {
	// No need to check connection, because we use library for Alexa
	return 200, nil
}

func (ag AlexaGrabber) GetServiceInfo() Service {
	return ag.Service
}

func (ag AlexaGrabber) Do(featuredHost string) (string, error) {
	globalRank, err := alexa.GlobalRank(featuredHost)
	if err != nil {
		fmt.Printf("Alexa.Do(): %s\n", err)
		return "No rank", nil
	}
	fmt.Printf("Alexa.Do(): %s rank in alexa is %s\n", featuredHost, globalRank)
	return globalRank, nil
}
