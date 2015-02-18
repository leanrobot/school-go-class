package netauth

import (
	"bitbucket.org/thopet/timeserver/config"
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var (
	httpAuthUrl string
	client      http.Client
)

func init() {
	httpAuthUrl = fmt.Sprintf("http://%s:%d", config.AuthUrl, config.AuthPort)
	log.Debug(httpAuthUrl)

	// setup the timeout transport for get200

	transport := http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr,
				time.Duration(config.AuthTimeout)*time.Millisecond)
		},
	}
	client = http.Client{
		Transport: &transport,
	}

	// test that the authserver is running.
	statusUrl := fmt.Sprintf(httpAuthUrl + "/status")
	// causes a panic if communication cannot be established with the
	// authserver.
	get200(statusUrl)
}

func Name(uuid string) (string, error) {
	url := fmt.Sprintf("%s/get?cookie=%s", httpAuthUrl, uuid)

	resp, err := get200(url)
	if err != nil {
		return "", err
	}
	name := getBodyAsString(resp.Body)
	if len(name) < 1 {
		return "", errors.New("No name returned")
	}
	return name, nil
}

func SetName(uuid string, name string) error {
	url := fmt.Sprintf("%s/set?cookie=%s&name=%s",
		httpAuthUrl, uuid, name)

	_, err := get200(url)
	if err != nil {
		return err
	}
	return nil
}

func ClearName(uuid string) error {
	url := fmt.Sprintf("%s/clear?cookie=%s", httpAuthUrl, uuid)

	if _, err := get200(url); err != nil {
		return err
	}
	return nil
}

// PRIVATE HELPERS ==========

func get200(url string) (res *http.Response, err error) {
	log.Debugf("making request to: %s", url)

	resp, err := client.Get(url)
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
