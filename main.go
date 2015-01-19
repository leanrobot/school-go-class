/*
Assignment 2: Personalized Time Server.
Written by Tom Petit (c) 2015
Winter 2015, CSS 490 - Tactical Software Engineering
*/

// The timeserver package contains a simple time server created for assignment one.
package main

import (
	"bitbucket.org/thopet/assign2/lib"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

/* ====== Learning Resources ======
- Standard Library Documentation
- https://golang.org/doc/articles/wiki/
- http://blog.golang.org/constants
- http://stackoverflow.com/questions/17891226/golang-operator-difference-between-vs
- https://httpd.apache.org/docs/1.3/logs.html#common
- https://ruin.io/2014/godoc-homebrew-go/
*/

/*
TODO:
	- add handler for view
		- / and /index.html
			- if logged in, write "Greetings, <name>".
			- if not, then show the login form.
		- /login/?name=...
			- if not logged in, generate a cookie with uuid. then redirect to /
			- no name param? say "C'mon, I need a name."
		- /logout/
			- clear the uuid cookie if it exists. Display goodbye for 10 seconds,
				then redirect to the home page.
		- /time/
			- if logged in, show "The time is now <time>, <name>."
			- if not, just show the time per assignment one.

*/

var userData *lib.UserData

//BUG URLs that are prefixed with '/time/' are still recognized as valid.
//For instance, '/time/notvalid' still returns the time and 200.

// Main method for the timeserver. Since this assignment is so small,
// all the setup and breakdown happens within this one function.
func main() {
	// setup and parse the arguments.
	var port *int = flag.Int("port", DEFAULT_PORT, "port to launch webserver on, default is 8080")
	var versionPrint *bool = flag.Bool("V", false, "Display version information")
	flag.Parse()

	if *versionPrint {
		fmt.Println(VERSION_INFO)
		os.Exit(0)
	}

	// setup and start the webserver.
	var portString string = fmt.Sprintf(":%d", *port)

	userData = lib.NewUserData()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/time/", timeHandler)

	fmt.Printf("Timeserver listening on 0.0.0.0%s\n", portString)
	err := http.ListenAndServe(portString, nil)

	if err != nil {
		log.Fatal("TimeServer Failure: ", err)
	}

	fmt.Println("Timeserver exiting..")
}

func getUsername(req *http.Request) (username string, err error) {

	if cookie, err := req.Cookie("assign2"); err == nil {
		uuid := lib.Uuid(cookie.Value)
		if username, err := userData.GetUser(uuid); err == nil {
			return username, nil
		}
	}
	return "", errors.New("No valid user uuid was found with an associated name")
}

func indexHandler(resStream http.ResponseWriter, req *http.Request) {
	// Extra handling is required to add 404 pages.
	if url := req.URL.Path; url != "/" || url != "/index.html" {
		notFoundHandler(resStream, req)
		return
	}

	// set the username based on the cookie if possible
	// usernameInsert := ""
	// uuid, err := getUuid(req)
	// if err == nil {
	// 	if username, err := userData.GetUser(uuid); err == nil {
	// 		usernameInsert = ", " + username
	// 	}
	// }
	username, err := getUsername(req)

	if err == nil {
		// a username was found, greet them.
		insertUsernameHtml := fmt.Sprintf(indexHtml.GetHtml(), username)
		io.WriteString(resStream, insertUsernameHtml)
	} else {
		// no username was found, display the login page.
		io.WriteString(resStream, loginHtml.GetHtml())
	}

}

func timeHandler(resStream http.ResponseWriter, req *http.Request) {
	var curTime string = time.Now().Local().Format(TIME_LAYOUT)
	io.WriteString(resStream, fmt.Sprintf(timeHtml.GetHtml(), curTime, ""))

	logRequest(req, http.StatusOK)
}

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

type HtmlData struct {
	head string
	body string
}

func (hd *HtmlData) GetHtml() string {
	return fmt.Sprintf(BASE_HTML, hd.head, hd.body)
}

const TIME_LAYOUT = "3:04:05 PM"
const DEFAULT_PORT = 8080

const VERSION_INFO = `
Simple Time Server. Written by Tom Petit.
Winter 2015, CSS 490 - Tactical Software Engineering

version: 1.0_assign1
`

const BASE_HTML = `
<html>
<head>%s</head>
<body>%s</body>
</html>
`

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

	indexHtml = HtmlData{
		head: "<style>#text { font-style: italic }</style>",
		body: "<p>Greetings, %s</p>",
	}

	logoutHtml = HtmlData{
		head: `<META http-equiv="refresh" content="10;URL=/">`,
		body: `<p>Good-bye.</p>`,
	}

	notFoundHtml = HtmlData{
		head: "",
		body: `<p>These are not the URLs you're looking for.</p>`,
	}
)
