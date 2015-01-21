/*
Assignment 2: Personalized Time Server.
Written by Tom Petit (c) 2015
Winter 2015, CSS 490 - Tactical Software Engineering

version: 2.0_assign2
*/

// The timeserver package contains a simple time server created for assignment one.
package main

import (
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

	http.HandleFunc("/", Handler404)
	http.HandleFunc("/time/", Handler200)

	fmt.Printf("Timeserver listening on 0.0.0.0%s\n", portString)
	err := http.ListenAndServe(portString, nil)

	if err != nil {
		log.Fatal("TimeServer Failure: ", err)
	}

	fmt.Println("Timeserver exiting..")
}

// Handler200 is the web handler for hits to /time on the webserver. This is
// considered the only valid url on the timeserver.
func Handler200(resStream http.ResponseWriter, req *http.Request) {
	var curTime string = time.Now().Local().Format(TIME_LAYOUT)
	var utcTime string = time.Now().UTC().Format(MILITARY_TIME_LAYOUT)
	io.WriteString(resStream, fmt.Sprintf(html_200, curTime, utcTime))

	logRequest(req, http.StatusOK)
}

// Handler404 is the handler for pages that do not exist. For this time server,
// that is everything except for "/time/"
func Handler404(resStream http.ResponseWriter, req *http.Request) {
	resStream.WriteHeader(http.StatusNotFound)
	io.WriteString(resStream, html_404)

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

const TIME_LAYOUT = "3:04:05 PM"
const MILITARY_TIME_LAYOUT = "15:04:05"

const DEFAULT_PORT = 8080

const VERSION_INFO = `
Simple Time Server. Written by Tom Petit.
Winter 2015, CSS 490 - Tactical Software Engineering

version: 1.0_assign1
`

// The html returned for a successful request. The time needs to be inserted
// as a string using Printf or Sprintf.
const html_200 = `
<html>
<head>
	<style>
		p { font-size: xx-large }
		span.time { color : red }
	</style>
</head>
<body>
	<p>The time is now <span class="time">%s</span> (%s UTC).</p>
</body>
</html>
`

// The html returned for an unknown page on the webserver.
const html_404 = `
<html>
<body>
	<p>These are not the URLs you're looking for.</p>
</body>
</html>
`
