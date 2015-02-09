package auth

import (
	"errors"
	"net/http"
	"io/ioutil"
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

// makes a request to the remote auth server to perform a login.
func (ca *CookieAuth) loginRequest(username string) (uuid Uuid, err error) {
	url := fmt.Sprintf("http://"+ca.authUrl+"/set?name=%s", username)
	log.Debugf("making request to: %s", url)
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			uuid := Uuid(body)
			return uuid, nil
		}
		return Uuid(""), errors.New("No valid Uuid returned.")
	}
	return Uuid(""), errors.New("login name not provided.")
}

/*
logout contains the logic for removing a uuid -> user association (via UserData)
struct and sets the cookie for removal.

BUG: the cookie is not deleted by the client. logout is still effective though.
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

// make a request to the remote auth server to perform a logout on that uuid.
func (ca *CookieAuth) logoutRequest(uuid Uuid) error {
	url := fmt.Sprintf("http://"+ca.authUrl+"/clear?uuid=%s", uuid)
	if resp, err := http.Get(url); err == nil {
		// check response is within 2xx range.
		if 200 <= resp.StatusCode && resp.StatusCode < 300 {
			return nil
		}
	}
	return errors.New("Logout Error.")
}

/*
getUsername retrieves the cookie from the request. It then uses the uuid
to retrieve the username from the UserData struct, returning the value.
*/
func (ca *CookieAuth) GetUsername(req *http.Request) (username string, err error) {
	cookie, err := req.Cookie(COOKIE_NAME)
	if err == nil {
		uuid := Uuid(cookie.Value)
		var username string
		username, err = ca.getRequest(uuid)
		if err == nil {
			return username, nil
		}
		return "", err
	}
	return "", err
}

func (ca *CookieAuth) getRequest(uuid Uuid) (username string, err error) {
	url := fmt.Sprintf("http://"+ca.authUrl+"/get?uuid=%s", uuid)
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			name := string(body)
			return name, nil
		}
	}
	return "", errors.New("Get User Error.")
}
