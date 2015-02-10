package main

import (
	"bitbucket.org/thopet/timeserver/auth"
	"bitbucket.org/thopet/timeserver/server"
	"bitbucket.org/thopet/timeserver/config"
	"net/http"
	log "github.com/cihub/seelog"
	"fmt"
	"io"
	"time"
)

var (
	users *auth.UserData 
)

func main() {

	// initialize the UserData manager.
	users = auth.NewUserData()

	// View Handler and patterns
	vh := server.NewStrictHandler()
	// TODO vh.NotFoundHandler
	vh.HandlePattern("/get", getName)
	vh.HandlePattern("/set", setName)
	vh.HandlePattern("/clear", clearName)

	portString := fmt.Sprintf(":%d", config.AuthPort)
	server := http.Server{
		Addr:    portString,
		Handler: vh,
	}

	log.Infof("authserver listening on port %d", config.AuthPort)
	err := server.ListenAndServe()

	if err != nil {
		log.Critical("authserver Failure: ", err)
	}

	log.Info("authserver exiting..")
}

func getName(res http.ResponseWriter, req *http.Request) {
	defer logRequest(req, http.StatusOK)
	uuid := auth.Uuid(req.FormValue("uuid"))
	if len(uuid) > 0 { // valid request path, return 200 and username
		name, err := users.GetUser(uuid)
		if err == nil {
			io.WriteString(res, name)
			return
		}
	}
	error400(res)
}

func setName(res http.ResponseWriter, req *http.Request) {
	defer logRequest(req, http.StatusOK)
	name := req.FormValue("name")
	if len(name) > 0 { // valid request path, return 200
		uuid, _ := users.AddUser(name)
		io.WriteString(res, uuid.String())
	} else { // non-valid request, return 400
		res.WriteHeader(http.StatusBadRequest)
	}
}

func clearName(res http.ResponseWriter, req *http.Request) {
	defer logRequest(req, http.StatusOK)
	uuid := auth.Uuid(req.FormValue("uuid"))
	if len(uuid) > 0 { // valid request path, return 200 and username
		ok := users.RemoveUser(uuid)
		if !ok {
			res.WriteHeader(http.StatusBadRequest)
		}
	} else { // non-valid request, return 400
		res.WriteHeader(http.StatusBadRequest)
	}
}

func logRequest(req *http.Request, statusCode int) {
	var requestTime string = time.Now().Format(time.RFC1123Z)

	fmt.Printf(`%s - [%s] "%s %s %s" %d -`+"\n",
		req.Host, requestTime, req.Method, req.URL.String(), req.Proto,
		statusCode)
}

func error400(res http.ResponseWriter) {
	log.Debug("Invalid Request [400]")
	res.WriteHeader(http.StatusBadRequest)
}
