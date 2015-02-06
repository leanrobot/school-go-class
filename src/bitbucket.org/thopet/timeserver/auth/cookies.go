package auth

import (
	"errors"
	"net/http"
	"io/ioutil"
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
	uuid, _ := ca.users.AddUser(username)

	// create a cookie
	cookie := http.Cookie{
		Name:   COOKIE_NAME,
		Path:   "/",
		Value:  string(uuid),
		MaxAge: 604800, // 7 days
	}
	http.SetCookie(res, &cookie)
	return nil
}

// makes a request to the remote auth server to perform a login.
func (ca *CookieAuth) loginRequest(username string) (uuid Uuid, err error) {
	url := "http://"+ca.authUrl+"/set?name=%s"
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			uuid := Uuid(body)
			return uuid, nil
		}
	}
	return Uuid(""), errors.New("Login Error.")
}

// make a request to the remote auth server to perform a logout on that uuid.
func logoutRequest(uuid Uuid) err error {
	url := "http://"+ca.authUrl+"/clear?name=%s"
	if resp, err := http.Get(url); err == nil {
		// check response is within 2xx range.
		if 200 <= resp.StatusCode && resp.StatusCode < 300 {
			return nil
		}
	}
	return Uuid(""), errors.New("Logout Error.")
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
		ca.users.RemoveUser(Uuid(cookie.Value))

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
	if cookie, err := req.Cookie(COOKIE_NAME); err == nil {
		uuid := Uuid(cookie.Value)
		if username, err := ca.users.GetUser(uuid); err == nil {
			return username, nil
		}
	}
	return "", errors.New("No valid user uuid was found with an associated name")
}
