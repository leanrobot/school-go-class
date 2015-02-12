package netauth

import (
	"fmt"
	"bitbucket.org/thopet/timeserver/config"
	"net/http"
	"io"
	"io/ioutil"
	log "github.com/cihub/seelog"
	"errors"
)

var (
	httpAuthUrl string
)

func init() {
	httpAuthUrl := "http://"+config.AuthUrl

	// test that the authserver is running.
	statusUrl := fmt.Sprintf(httpAuthUrl+"/status")
	// causes a panic if communication cannot be established with the
	// authserver.
	get200(statusUrl) 
}

func Name(uuid string) (string, error) {
	url := fmt.Sprintf(httpAuthUrl+"/get?cookie=%s", uuid)

	resp, err := get200(url)
	if err != nil {
		return "", err
	}
	name := getBodyAsString(resp.Body)
	return name, nil
}

func SetName(uuid string, name string) error {
	url := fmt.Sprintf(httpAuthUrl+"/set?cookie=%s&name=%s", uuid, name)
	
	_, err := get200(url)
	if err != nil {
		return err
	}
	return nil
}

func ClearName(uuid string) error {
	url := fmt.Sprintf(httpAuthUrl+"/clear?cookie=%s", uuid)

	if _, err := get200(url); err != nil {
		return err
	}
	return nil
}

// PRIVATE HELPERS ==========

func get200(url string) (res *http.Response, err error) {
	log.Debugf("making request to: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	status := resp.StatusCode
	if 200 <= status && status < 300 {
		return resp, nil
	}
	return nil, errors.New("Not a 2xx response.")
}

func getBodyAsString(body io.ReadCloser) string {
	defer body.Close()
	if body, err := ioutil.ReadAll(body); err == nil {
		contents := string(body)
		return contents
	}
	panic("Could not read body into string")
}