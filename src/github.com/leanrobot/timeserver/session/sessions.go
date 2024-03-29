package session

import (
	"crypto/md5"
	"fmt"
	"github.com/leanrobot/timeserver/config"
	"github.com/leanrobot/timeserver/cookie"
	"github.com/leanrobot/timeserver/netauth"
	"io"
	"net/http"
	"time"
)

var (
	sessionName string
)

func init() {
	sessionName = config.SESSION_NAME
}

// TODO error return value?
func Create(res http.ResponseWriter, name string) error {
	uuid := uuidGen()
	err := netauth.SetName(uuid, name)
	if err != nil {
		return err
	}
	cookie.Create(res, sessionName, uuid)
	return nil
}

func Destroy(req *http.Request, res http.ResponseWriter) error {
	// if the cookie doesn't exist, session doesn't exist.
	uuid, err := cookie.Get(req, sessionName)
	// a cookie was not found for the session, so no need to delete.
	if err != nil {
		return nil
	}
	// a session exists that needs to be deleted.
	cookie.Clear(res, sessionName)

	err = netauth.ClearName(uuid)
	if err != nil {
		return err
	}
	return nil
}

func Username(req *http.Request) (string, error) {
	uuid, err := cookie.Get(req, sessionName)
	if err != nil {
		return "", err
	}

	name, err := netauth.Name(uuid)
	if err != nil {
		return "", err
	}
	return name, nil
}

func uuidGen() string {
	idHash := md5.New()
	io.WriteString(idHash, time.Now().Format(time.UnixDate))
	io.WriteString(idHash, time.Now().Format(time.UnixDate))

	uuid := fmt.Sprintf("%x", idHash.Sum(nil))
	return uuid
}
