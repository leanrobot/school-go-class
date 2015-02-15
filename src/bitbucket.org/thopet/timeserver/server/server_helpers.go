package server

import (
	"net/http"
	"bitbucket.org/thopet/timeserver/config"
	"time"
	log "github.com/cihub/seelog"
)

var (
	max int
	featureOn bool
	// the number of booleans in this queue represents the number of requests
	// which may run concurrently.
	queue chan bool
)

func init() {
	max = config.RequestLimit
	featureOn = true

	if max == 0 {
		featureOn = false
		log.Info("No Request Limit Enabled")
	}

	queue = make(chan bool, max)
	for i:=0; i<max; i++ {
		queue <- true
	}
}

func LimitRequests(h http.HandlerFunc) http.HandlerFunc {
	// if the feature isn't enable (max=0) then don't use closure.
	if !featureOn {
		return h
	}
	return func(res http.ResponseWriter, req *http.Request) {
		queue := queue
		select {
		case <- queue:
			h(res, req)
			queue <- true
		default:
			Error502(res, req)
		}
	}


}

func Error400(res http.ResponseWriter, req *http.Request) {
	LogRequest(req, http.StatusBadRequest)
	res.WriteHeader(http.StatusBadRequest)
}

func Error502(res http.ResponseWriter, req *http.Request) {
	LogRequest(req, http.StatusServiceUnavailable)
	res.WriteHeader(http.StatusServiceUnavailable)
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
func LogRequest(req *http.Request, statusCode int) {
	var requestTime string = time.Now().Format(time.RFC1123Z)

	log.Infof(`%s - [%s] "%s %s %s" %d -`,
		req.Host, requestTime, req.Method, req.URL.String(), req.Proto,
		statusCode)
}