package grabber

import (
	"net/http"
	"fmt"
	"io/ioutil"
)

type Service struct {
	Name string
	URL  string
}

func GeneralCheckConnection(s Service) (int, error) {
	resp, err := http.Get(s.URL)
	if err != nil {
		fmt.Println(err)
		return resp.StatusCode, err
	}
	defer resp.Body.Close()
	return 200, nil
}

func GeneralDo(s Service) (string, error) {
	resp, err := http.Get(s.URL)
	if err != nil {
		fmt.Println(err)
		return resp.Status, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(body), nil
}

