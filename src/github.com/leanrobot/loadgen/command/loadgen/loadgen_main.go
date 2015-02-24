package main

import (
	"flag"
	"fmt"
	"github.com/leanrobot/loadgen/counter"
	"net"
	"net/http"
	"time"
)

/*
TODO
- count the errors
- wait an additional timeout after finishing main loop to make sure
		that requests have finished.
*/

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

func initConfig() {
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

func main() {
	initConfig()

	// setup a client with the timeout.
	transport := http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr,
				time.Duration(Timeout)*time.Millisecond)
		},
	}
	client = http.Client{
		Transport: &transport,
	}

	// run the launcher process.
	go launcher()
	time.Sleep(Runtime)

	// report the statistics
	stats := counter.Export()
	keys := []int{
		100, 200, 300, 400, 500,
	}

	for _, key := range keys {
		fmt.Printf("%ds\t\t\t%d\n", key, stats[string(key)])
	}

	fmt.Println("actualy map:")
	for key, value := range stats {
		fmt.Printf("%s\t\t\t%d\n", key, value)
	}

}

func launcher() {
	interval := time.Duration(Burst*1000000/Rate) * time.Microsecond
	ticker := time.Tick(interval)

	for {
		<-ticker
		for i := 0; i < Burst; i++ {
			go worker()
		}
	}
}

func worker() {
	resp, err := client.Get(Url)
	if err != nil {
		panic(err)
	}

	// get the status code
	status := resp.StatusCode
	statusKey := fmt.Sprintf("%d", status)

	// increment the status variable.
	counter.Increment(statusKey)
}
