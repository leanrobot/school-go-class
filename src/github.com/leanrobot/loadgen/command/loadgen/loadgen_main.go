package main

import (
	"flag"
	"fmt"
	"github.com/leanrobot/loadgen/counter"
	"net/http"
	"time"
)

var (
	Rate    int
	Burst   int
	Timeout time.Duration
	Runtime time.Duration
	Url     string

	client http.Client
)

const (
	DEFAULT_RATE    int           = 200
	DEFAULT_BURST   int           = 20
	DEFAULT_TIMEOUT time.Duration = 1000 * time.Millisecond
	DEFAULT_RUNTIME time.Duration = 10 * time.Second
	DEFAULT_URL     string        = "http://localhost:8080/time"
)

func initFlags() {
	flag.IntVar(&Rate, "rate", DEFAULT_RATE,
		"the number of requests to send per second")
	flag.IntVar(&Burst, "burst", DEFAULT_BURST,
		"The number of concurrent requests that will be issues")
	flag.DurationVar(&Timeout, "timeout", DEFAULT_TIMEOUT,
		"The timeout when issuing requests")
	flag.DurationVar(&Runtime, "runtime", DEFAULT_RUNTIME,
		"The duration to perform the test for")
	flag.StringVar(&Url, "url", DEFAULT_URL, "The url to test.")

	flag.Parse()
}

func initClient() {
	client = http.Client{
		Timeout: Timeout,
	}
}

func main() {
	initFlags()
	// setup a client with the timeout.
	initClient()

	// run the launcher process.
	controller := make(chan bool)
	go launcher(controller)
	time.Sleep(Runtime)

	// kill the launcher and wait for requests to timeout.
	controller <- true
	time.Sleep(Timeout * 2)

	// report the statistics
	stats := counter.Export()
	keys := []string{
		"total", "100s", "200s", "300s", "400s", "500s", "error",
	}

	for _, key := range keys {
		fmt.Printf("%s:\t%d\n", key, stats[key])
	}
}

// launcher is a goroutine which launches the requests off to the remote page
// that is being tested. the controller channel is used to terminate this
// routine. when anything is received on this channel, the launcher closes the
// channel and terminates itself.
func launcher(controller chan bool) {
	interval := time.Duration(Burst*1000000/Rate) * time.Microsecond
	ticker := time.Tick(interval)

	for {
		select {
		case <-ticker:
			break
		case <-controller:
			close(controller)
			return
		}

		for i := 0; i < Burst; i++ {
			go worker()
		}
	}
}

/*
worker is a simple goroutine which makes a request to a remote page. it then
takes the status code from the response and counts them up by century using
the counter.
*/
func worker() {
	counter.Increment("total")

	resp, err := client.Get(Url)
	if err != nil {
		counter.Increment("error")
		return
	}
	defer resp.Body.Close()

	// get the status code by century.
	status := (resp.StatusCode / 100) * 100
	statusKey := fmt.Sprintf("%ds", status)

	// increment the status variable.
	counter.Increment(statusKey)
}
