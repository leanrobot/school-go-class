package cookie

import (
	"net/http"
)

func Create(res http.ResponseWriter, key string, value string) {
	// create a cookie
	cookie := newCookie(key, value)
	http.SetCookie(res, cookie)
}

func Get(req *http.Request, key string) (string, error) {
	cookie, err := req.Cookie(key)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// dumb clear, does not check for existence of cookie.
func Clear(res http.ResponseWriter, key string) {
	cookie := newCookie(key, "CLEARED_VALUE")
	cookie.MaxAge = -1
	http.SetCookie(res, cookie)
}

func newCookie(key string, value string) *http.Cookie{
	cookie := http.Cookie{
		Name:   key,
		Path:   "/",
		Value:  value,
		MaxAge: 604800, // 7 days
	}
	return &cookie
}