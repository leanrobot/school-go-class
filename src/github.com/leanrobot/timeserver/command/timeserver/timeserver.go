/*
Assignment 2: Personalized Time Server.
Written by Tom Petit (c) 2015
Winter 2015, CSS 490 - Tactical Software Engineering
*/

// The timeserver package contains a simple time server created for assignment one.
package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/leanrobot/timeserver/config"
	"github.com/leanrobot/timeserver/server"
	"github.com/leanrobot/timeserver/session"
	"html/template"
	"math/rand"
	"net/http"
	"time"
)

const (
	TIME_LAYOUT          = "3:04:05 PM"
	MILITARY_TIME_LAYOUT = "15:04:05"
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
	// setup and start the webserver.
	portString := fmt.Sprintf(":%d", config.Port)

	initTemplates(config.TemplatesDir)

	// custom handler with strict url pattern matching
	vh := server.NewStrictHandler()
	vh.NotFoundHandler = notFoundHandler
	vh.HandlePatterns([]string{"/", "/index.html"}, indexHandler)
	vh.HandlePattern("/time/", server.LimitRequests(timeHandler))
	vh.HandlePattern("/login/", loginHandler)
	vh.HandlePattern("/logout/", logoutHandler)
	vh.HandlePattern("/about/", aboutHandler)
	vh.ServeStaticFile("/css/style.css", config.TemplatesDir+"/style.css")

	log.Infof("Timeserver listening on 0.0.0.0%s", portString)
	err := http.ListenAndServe(portString, vh)

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
	defer server.LogRequest(req, http.StatusOK)

	username, err := session.Username(req)
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
	defer server.LogRequest(req, http.StatusFound)

	// get the requested username
	username := req.FormValue("name")
	if len(username) < 1 {
		renderBaseTemplate(res, "login_error.html", nil)
	} else {
		err := session.Create(res, username)
		if err != nil {
			log.Error(err)
		}
		http.Redirect(res, req, "/index.html", http.StatusFound)
	}
}

// logoutHandler is the view for the logout resource.
func logoutHandler(res http.ResponseWriter, req *http.Request) {
	defer server.LogRequest(req, http.StatusFound)

	session.Destroy(req, res)
	renderBaseTemplate(res, "logout.html", nil)
}

// timeHandler is the view for the time resource.
func timeHandler(res http.ResponseWriter, req *http.Request) {
	defer server.LogRequest(req, http.StatusOK)

	// replace empty string with the username text if logged in.
	data := struct {
		Time         string
		MilitaryTime string
		Username     string
	}{
		Time:         time.Now().Local().Format(TIME_LAYOUT),
		MilitaryTime: time.Now().UTC().Format(MILITARY_TIME_LAYOUT),
	}

	// TODO implement random load simulation
	wait := rand.NormFloat64()*float64(config.Deviation) +
		float64(config.AvgResponse)
	if wait < 0 {
		wait = 0
	}
	log.Infof("sleep duration is %d", wait)
	time.Sleep(time.Duration(wait) * time.Millisecond)

	if username, err := session.Username(req); err == nil {
		data.Username = username
	}

	renderBaseTemplate(res, "time.html", data)
}

func aboutHandler(res http.ResponseWriter, req *http.Request) {
	defer server.LogRequest(req, http.StatusOK)

	res.WriteHeader(http.StatusOK)
	renderBaseTemplate(res, "about_us.html", nil)
}

// notFoundHandler is the view for the global 404 resource.
func notFoundHandler(res http.ResponseWriter, req *http.Request) {
	defer server.LogRequest(req, http.StatusNotFound)

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
