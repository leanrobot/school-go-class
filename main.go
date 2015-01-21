/*
Assignment 2: Personalized Time Server.
Written by Tom Petit (c) 2015
Winter 2015, CSS 490 - Tactical Software Engineering
*/

// The timeserver package contains a simple time server created for assignment one.
package main


import (
	"bitbucket.org/thopet/timeserver/lib"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"sync"
)

const VERSION = "assignment-02.rc02"

/*
TODO:
	- Refactor main.go into modules.
	- BUG: login can happen during an existing login
		causing an orphaned UUID in the data.
	- /logout/
		- BUG(assign2): client does not delete cookie (for some reason)
			- Still is invalidated (so good).
	- TODO(assign2): error handling for hash collision?

Learning in this assignment:
	- structs and receiver functions for structs.
	- slices, maps.
	- Custom server handlers
	- Error handling in go.
	- Multiple return values
	- First-order functions.
	- Cookies
*/

// userData holds all the user information with a uuid association.
var userData *lib.UserData
// dataMutex is used to lock the shared userData struct.
var dataMutex *sync.Mutex

// Main method for the timeserver.
func main() {
	// setup and parse the arguments.
	var port *int = flag.Int("port", DEFAULT_PORT, "port to launch webserver on, default is 8080")
	var versionPrint *bool = flag.Bool("V", false, "Display version information")
	flag.Parse()

	if *versionPrint {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	// Initialize global data.
	userData = lib.NewUserData()
	dataMutex = new(sync.Mutex)

	// setup and start the webserver.
	var portString string = fmt.Sprintf(":%d", *port)

	// View Handler and patterns
	vh := lib.NewStrictHandler()
	vh.NotFoundHandler = notFoundHandler
	vh.HandlePatterns([]string{"/", "/index.html"}, indexHandler)
	vh.HandlePattern("/time/", timeHandler)
	vh.HandlePattern("/login/", loginHandler)
	vh.HandlePattern("/logout/", logoutHandler)

	server := http.Server{
		Addr:    portString,
		Handler: vh,
	}

	fmt.Printf("Timeserver listening on 0.0.0.0%s\n", portString)
	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("TimeServer Failure: ", err)
	}

	fmt.Println("Timeserver exiting..")
}

/* 
getUsername retrieves the cookie from the request. It then uses the uuid
to retrieve the username from the UserData struct, returning the value.
*/
func getUsername(req *http.Request) (username string, err error) {
	if cookie, err := req.Cookie(COOKIE_NAME); err == nil {
		uuid := lib.Uuid(cookie.Value)
		if username, err := userData.GetUser(uuid); err == nil {
			return username, nil
		}
	}
	return "", errors.New("No valid user uuid was found with an associated name")
}

// indexHandler is the view for the index resource.
func indexHandler(resStream http.ResponseWriter, req *http.Request) {
	username, err := getUsername(req)
	if err == nil {
		// a username was found, greet them.
		insertUsernameHtml := fmt.Sprintf(indexHtml.GetHtml(), username)
		io.WriteString(resStream, insertUsernameHtml)
	} else {
		// no username was found, display the login page.
		io.WriteString(resStream, loginHtml.GetHtml())
	}

	logRequest(req, http.StatusOK)
}

/*
login create a new uuid -> username association using the UserData struct.
A cookie is then created to store the uuid.
*/ 
func login(res http.ResponseWriter, username string) error {
	// TODO: error handling for hash collision?

	dataMutex.Lock()
	fmt.Fprintln(os.Stderr, "login(): Mutex Lock")

	uuid, _ := userData.AddUser(username)

	dataMutex.Unlock()
	fmt.Fprintln(os.Stderr, "login(): Mutex Unlock")

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

// loginHandler is the view for the login resource.
func loginHandler(res http.ResponseWriter, req *http.Request) {
	// get the requested username
	username := req.FormValue("name")
	if len(username) < 1 {
		io.WriteString(res, loginHtmlError.GetHtml())
	} else {
		login(res, username)
		http.Redirect(res, req, "/index.html", http.StatusFound)
	}
	
	logRequest(req, http.StatusFound)
}

/*
logout contains the logic for removing a uuid -> user association (via UserData)
struct and sets the cookie for removal.

BUG: the cookie is not deleted by the client. logout is still effective though.
TODO: add error handling?
*/
func logout(res http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie(COOKIE_NAME)

	fmt.Fprintln(os.Stderr, "logout(): Mutex Lock")
	dataMutex.Lock()
	userData.RemoveUser(lib.Uuid(cookie.Value))
	dataMutex.Unlock()
	fmt.Fprintln(os.Stderr, "logout(): Mutex Unlock")

	cookie.MaxAge = -1
	cookie.Value = "LOGGED_OUT_CLEAR_DATA"
	http.SetCookie(res, cookie)
}

// logoutHandler is the view for the logout resource.
func logoutHandler(res http.ResponseWriter, req *http.Request) {
	logout(res, req)
	io.WriteString(res, logoutHtml.GetHtml())

	logRequest(req, http.StatusFound)
}

// timeHandler is the view for the time resource.
func timeHandler(resStream http.ResponseWriter, req *http.Request) {
	// replace empty string with the username text if logged in.
	usernameInsert := ""
	if username, err := getUsername(req); err == nil {
		usernameInsert = ", " + username
	}

	var curTime string = time.Now().Local().Format(TIME_LAYOUT)
	io.WriteString(resStream, fmt.Sprintf(timeHtml.GetHtml(),
		curTime, usernameInsert))

	logRequest(req, http.StatusOK)
}

// notFoundHandler is the view for the global 404 resource.
func notFoundHandler(resStream http.ResponseWriter, req *http.Request) {
	resStream.WriteHeader(http.StatusNotFound)
	io.WriteString(resStream, notFoundHtml.GetHtml())

	logRequest(req, http.StatusNotFound)
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

	fmt.Printf(`%s - [%s] "%s %s %s" %d -`,
		req.Host, requestTime, req.Method, req.URL.String(), req.Proto,
		statusCode)
	fmt.Println("")
}

/*
HtmlData is a simple struct for holding simple template info.
*/
type HtmlData struct {
	head string
	body string
}

// GetHtml renders the head and body of an HtmlData into the BASE_HTML template.
func (hd *HtmlData) GetHtml() string {
	return fmt.Sprintf(BASE_HTML, hd.head, hd.body)
}

const COOKIE_NAME = "timeserver_assignment2_tompetit"
const TIME_LAYOUT = "3:04:05 PM"
const DEFAULT_PORT = 8080

const BASE_HTML = `
<html>
<head>%s</head>
<body>%s</body>
</html>
`

type HandlerFunc func(http.ResponseWriter, *http.Request)

// TODO(assign2): constant structs?
// These ARE constant. They are not declared as const because they are structs.
var (
	timeHtml = HtmlData{
		head: `
			<style>
				p { font-size: xx-large }
				span.time { color : red }
			</style>
		`,
		body: `<p>The time is now <span class="time">%s</span>%s.</p>`,
	}

	loginHtml = HtmlData{
		head: "",
		body: `
			<p>
				<form method="POST" action="/login/">
					What is your name, Earthling?
					<input type="text" name="name" size="50">
					<input type="submit">
				</form>
			</p>
		`,
	}

	loginHtmlError = HtmlData{
		body: `<p>C'mon, I need a name.</p>`,
	}

	indexHtml = HtmlData{
		head: "<style>#text { font-style: italic }</style>",
		body: "<p>Greetings, %s</p>",
	}

	logoutHtml = HtmlData{
		head: `<META http-equiv="refresh" content="3;URL=/">`,
		body: `<p>Good-bye.</p>`,
	}

	notFoundHtml = HtmlData{
		head: "",
		body: `<p>These are not the URLs you're looking for.</p>`,
	}
)
