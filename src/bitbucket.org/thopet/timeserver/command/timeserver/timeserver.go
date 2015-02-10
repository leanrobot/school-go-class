/*
Assignment 2: Personalized Time Server.
Written by Tom Petit (c) 2015
Winter 2015, CSS 490 - Tactical Software Engineering
*/

// The timeserver package contains a simple time server created for assignment one.
package main

import (
	"bitbucket.org/thopet/timeserver/auth"
	"bitbucket.org/thopet/timeserver/server"
	"bitbucket.org/thopet/timeserver/config"
	"fmt"
	log "github.com/cihub/seelog"
	"html/template"
	"net/http"
	"time"
)

const (
	TIME_LAYOUT          = "3:04:05 PM"
	MILITARY_TIME_LAYOUT = "15:04:05"
)

/*
TODO:
	- BUG: login can happen during an existing login
		causing an orphaned UUID in the data.
	- /logout/
		- BUG(assign2): client does not delete cookie (for some reason)
			- Still is invalidated (so good).
	- TODO(assign2): error handling for hash collision?
*/

var (
	// auth holds all the user information with a uuid association.
	cAuth *auth.CookieAuth
)


var templates = map[string]*template.Template{
	"index.html":       nil,
	"time.html":        nil,
	"login.html":       nil,
	"login_error.html": nil,
	"logout.html":      nil,
	"404.html":         nil,
	"about_us.html":    nil,
}

// Main method for the timeserver.
func main() {
	// initialize the auth module
	authHost := fmt.Sprintf("%s:%d", config.AuthUrl, config.AuthPort)
	cAuth = auth.NewCookieAuth(authHost)

	// setup and start the webserver.
	portString := fmt.Sprintf(":%d", config.Port)

	initTemplates(config.TemplatesDir)

	// View Handler and patterns
	vh := server.NewStrictHandler()
	vh.NotFoundHandler = notFoundHandler
	vh.HandlePatterns([]string{"/", "/index.html"}, indexHandler)
	vh.HandlePattern("/time/", timeHandler)
	vh.HandlePattern("/login/", loginHandler)
	vh.HandlePattern("/logout/", logoutHandler)
	vh.HandlePattern("/about/", aboutHandler)
	vh.ServeStaticFile("/css/style.css", config.TemplatesDir+"/style.css")

	server := http.Server{
		Addr:    portString,
		Handler: vh,
	}

	log.Infof("Timeserver listening on 0.0.0.0%s", portString)
	err := server.ListenAndServe()

	if err != nil {
		log.Critical("TimeServer Failure: ", err)
	}

	log.Info("Timeserver exiting..")
}

func initTemplates(templateDir string) {
	for key, _ := range templates {
		templatePath := func(filename string) string {
			return templateDir + "/" + filename
		}
		templates[key] = template.Must(template.ParseFiles(
			templatePath("base.html"),
			templatePath("menu.html"),
			templatePath(key),
		))
	}
}

// indexHandler is the view for the index resource.
func indexHandler(res http.ResponseWriter, req *http.Request) {
	defer logRequest(req, http.StatusOK)

	username, err := cAuth.GetUsername(req)
	if err == nil {
		data := struct{ Username string }{Username: username}
		// a username was found, greet them.
		renderBaseTemplate(res, "index.html", data)
	} else {
		log.Error(err)
		renderBaseTemplate(res, "login.html", nil)
	}
}

// loginHandler is the view for the login resource.
func loginHandler(res http.ResponseWriter, req *http.Request) {
	defer logRequest(req, http.StatusFound)

	// get the requested username
	username := req.FormValue("name")
	if len(username) < 1 {
		renderBaseTemplate(res, "login_error.html", nil)
	} else {
		err := cAuth.Login(res, username)
		if err != nil { log.Error(err) }
		http.Redirect(res, req, "/index.html", http.StatusFound)
	}
}

// logoutHandler is the view for the logout resource.
func logoutHandler(res http.ResponseWriter, req *http.Request) {
	defer logRequest(req, http.StatusFound)

	cAuth.Logout(res, req)
	renderBaseTemplate(res, "logout.html", nil)
}

// timeHandler is the view for the time resource.
func timeHandler(res http.ResponseWriter, req *http.Request) {
	defer logRequest(req, http.StatusOK)

	// replace empty string with the username text if logged in.
	data := struct {
		Time         string
		MilitaryTime string
		Username     string
	}{}

	if username, err := cAuth.GetUsername(req); err == nil {
		data.Username = username
	}
	data.Time = time.Now().Local().Format(TIME_LAYOUT)
	data.MilitaryTime = time.Now().UTC().Format(MILITARY_TIME_LAYOUT)

	renderBaseTemplate(res, "time.html", data)
}

func aboutHandler(res http.ResponseWriter, req *http.Request) {
	defer logRequest(req, http.StatusOK)

	res.WriteHeader(http.StatusOK)
	renderBaseTemplate(res, "about_us.html", nil)
}

// notFoundHandler is the view for the global 404 resource.
func notFoundHandler(res http.ResponseWriter, req *http.Request) {
	defer logRequest(req, http.StatusNotFound)

	res.WriteHeader(http.StatusNotFound)
	renderBaseTemplate(res, "404.html", nil)
}

func renderBaseTemplate(res http.ResponseWriter, templateName string, data interface{}) {
	var err error
	tmpl, ok := templates[templateName]
	if !ok {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	} else {
		err = tmpl.ExecuteTemplate(res, "base", data)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

// logRequest logs request data to stdout. The format conforms closely to
// Apache Common Log Format.
//
// Example:
// 127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326
//
// {host} {user} [{time}] "{method} {url} {protocol}/{version}" {response status code} {response size}
// The user and response size are not supported, a - is used to fill the space.
//
// Reference: https://httpd.apache.org/docs/1.3/logs.html#common
func logRequest(req *http.Request, statusCode int) {
	var requestTime string = time.Now().Format(time.RFC1123Z)

	log.Infof(`%s - [%s] "%s %s %s" %d -`,
		req.Host, requestTime, req.Method, req.URL.String(), req.Proto,
		statusCode)
}
