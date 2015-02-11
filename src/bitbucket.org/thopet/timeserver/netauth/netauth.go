package netauth

import (
	"bitbucket.org/thopet/timeserver/auth"
	"bitbucket.org/thopet/timeserver/config"
)

var (
	httpAuthUrl string
)

func init() {
	httpAuthUrl := "http://"+config.authUrl

	// test that the authserver is running.
	statusUrl := fmt.Sprintf(httpAuthUrl+"/status")
	// causes a panic if communication cannot be established with the
	// authserver.
	_, err = get200(statusUrl) 
	
}

func SetName(name string) (string, error) {
	uuid := uuidGen()

	url := fmt.Sprintf(httpAuthUrl+"/set?cookie=%s&name=%s", uuid, name)
	resp, err = get200(url)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func Name(uuid string) string {
	url := fmt.Sprintf(httpAuthUrl+"/get?cookie=%s", uuid)

	resp, err := get200(url)
	if err != nil {
		return "", err
	}
	name := getBodyAsString(resp.Body)
	return name, nil
}

func ClearName(uuid string) error {
	url := fmt.Sprintf(httpAuthUrl+"/clear?cookie=%s", uuid)

	if _, err := get200(url); err != nil {
		return err
	}
	return nil
}

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

// TODO
func uuidGen() string {
	return "TEST TEST TEST TEST"
}