package server

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/leanrobot/counter"
	"github.com/leanrobot/timeserver/config"
	"net/http"
	"time"
)

var (
	max       int
	featureOn bool
	// A semaphore channel. When the channel is drained of bool's
	// the semaphore is locked.
	sem chan bool
)

func init() {
	max = config.RequestLimit
	featureOn = true

	if max == 0 {
		featureOn = false
		log.Info("No Request Limit Enabled")
	}

	sem = make(chan bool, max)
	for i := 0; i < max; i++ {
		sem <- true
	}
}

func LimitRequests(h http.HandlerFunc) http.HandlerFunc {
	// if the feature isn't enable (max=0) then don't use closure.
	if !featureOn {
		return h
	}
	return func(res http.ResponseWriter, req *http.Request) {
		sem := sem
		select {
		case <-sem:
			h(res, req)
			sem <- true
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

	// log the century of the status code
	century := (statusCode / 100) * 100
	if century == 200 || century == 400 {
		counter.Increment(fmt.Sprintf("%ds", century))
	}

	// log 404s
	if statusCode == 404 {
		counter.Increment(fmt.Sprintf("%ds", statusCode))
	}

	log.Infof(`%s - [%s] "%s %s %s" %d -`,
		req.Host, requestTime, req.Method, req.URL.String(), req.Proto,
		statusCode)
}

/*
MonitorHandler is a http.HandlerFunc type which uses the counter package
in order to "export" the contents of the program counter to an external service,
AKA the monitoring service for assignment 6. The data is displayed in JSON.
*/
func MonitorHandler(res http.ResponseWriter, req *http.Request) {
	data := counter.Export()
	dataJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(dataJson)
}
