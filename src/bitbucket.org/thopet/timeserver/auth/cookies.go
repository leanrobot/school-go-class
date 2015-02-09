package auth

import (
	"errors"
	"net/http"
	"io/ioutil"
	"io"
	"fmt"
	log "github.com/cihub/seelog"
)

const (
	COOKIE_NAME = "timeserver_assignment3_tompetit"
)

type CookieAuth struct {
	users      *UserData
	cookieName string
	authUrl  string
}

func NewCookieAuth(authUrl string) *CookieAuth {
	return &CookieAuth {
		users:      NewUserData(),
		cookieName: COOKIE_NAME,
		authUrl: authUrl,
	}
}

/*
login create a new uuid -> username association using the UserData struct.
A cookie is then created to store the uuid.
*/
func (ca *CookieAuth) Login(res http.ResponseWriter, username string) error {
	// TODO: error handling for hash collision?
	uuid, err := ca.loginRequest(username)

	if err == nil {
		// create a cookie
		cookie := http.Cookie{
			Name:   COOKIE_NAME,
			Path:   "/",
			Value:  string(uuid),
			MaxAge: 604800, // 7 days
		}
		http.SetCookie(res, &cookie)
	}
	return err 
}

/*
logout contains the logic for removing a uuid -> user association (via UserData)
struct and sets the cookie for removal.

TODO: add error handling?
*/
func (ca *CookieAuth) Logout(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(COOKIE_NAME)

	// only perform logouts for users with a cookie.
	if err == nil {
		ca.logoutRequest(Uuid(cookie.Value))

		cookie.MaxAge = -1
		cookie.Value = "LOGGED_OUT_CLEAR_DATA"
		http.SetCookie(res, cookie)
	}
}

/*
getUsername retrieves the cookie from the request. It then uses the uuid
to retrieve the username from the UserData struct, returning the value.
*/
func (ca *CookieAuth) GetUsername(req *http.Request) (username string, err error) {
	cookie, err := req.Cookie(COOKIE_NAME)
	if err == nil {
		uuid := Uuid(cookie.Value)
		username, err := ca.getRequest(uuid)
		if err == nil {
			log.Debugf("username found: [%s], %s", username, err)
			return username, nil
		}
		return "", err
	}
	return "", err
}

// PRIVATE HELPERS =================

// makes a request to the remote auth server to perform a login.
func (ca *CookieAuth) loginRequest(username string) (uuid Uuid, err error) {
	url := fmt.Sprintf("http://"+ca.authUrl+"/set?name=%s", username)

	resp, err := get200(url)
	if err != nil {
		return Uuid(""), err
	}
	uuid = Uuid(getBodyAsString(resp.Body))
	return uuid, nil
}

// make a request to the remote auth server to perform a logout on that uuid.
func (ca *CookieAuth) logoutRequest(uuid Uuid) error {
	url := fmt.Sprintf("http://"+ca.authUrl+"/clear?uuid=%s", uuid)

	if _, err := get200(url); err != nil {
		return err
	}
	return nil
}

func (ca *CookieAuth) getRequest(uuid Uuid) (username string, err error) {
	url := fmt.Sprintf("http://"+ca.authUrl+"/get?uuid=%s", uuid)

	resp, err := get200(url)
	if err != nil {
		return "", err
	}
	name := getBodyAsString(resp.Body)
	return name, nil
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

func getBodyAsString(body io.ReadCloser) string {
	defer body.Close()
	if body, err := ioutil.ReadAll(body); err == nil {
		contents := string(body)
		return contents
	}
	panic("Could not read body into string")
}
