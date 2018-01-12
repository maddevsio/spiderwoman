package grabber

import (
	"fmt"

	"github.com/domainr/whois"
)

type WhoisGrabber struct {
	Service
}

func (ag WhoisGrabber) CheckConnection() (int, error) {
	return 200, nil
}

func (ag WhoisGrabber) GetServiceInfo() Service {
	return ag.Service
}

func (ag WhoisGrabber) Do(featuredHost string) (string, error) {
	request, err := whois.NewRequest(featuredHost)
	if err != nil {
		fmt.Printf("Whois.Do() NewRequest: %s\n", err)
		return err.Error(), nil // return error string, need to save it
	}
	response, err := whois.DefaultClient.Fetch(request)
	if err != nil {
		fmt.Printf("Whois.Do() DefaultClient: %s\n", err)
		return err.Error(), nil // return error string, need to save it
	}

	body := string(response.Body)
	fmt.Println(body)

	return body, nil
}
