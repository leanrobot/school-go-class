package loadgen

import (
	"flag"
)

type Config struct {
	RequestRate time.Duration
	Burst       int
	Timeout     time.Duration
	Runtime     time.Duration
	Url         string
}

const (
	DEFAULT_REQUEST_RATE time.Duration = 200 * time.Second
	DEFAULT_BURST        int           = 20
	DEFAULT_TIMEOUT      time.Duration = 1000 * time.Millisecond
	DEFAULT_RUNTIME      time.Duration = 10 * time.Second
	DEFAULT_URL          string        = "http://localhost:8080/time"
)

var (
	config *Config
)

func initConfig() *Config {
	config = new(Config)
	flag.DurationVar(&config.RequestRate, "rate")
	flag.IntVar(&config.Burst, "burst", DEFAULT_BURST,
		"The number of concurrent requests that will be issues")
	flag.DurationVar(&config.Timeout, "timeout", DEFAULT_TIMEOUT,
		"The timeout when issuing requests")
	flag.DurationVar(&config.Runtime, "runtime", DEFAULT_RUNTIME,
		"The duration to perform the test for")
	flag.StringVar(&config.Url, "url", DEFAULT_URL, "The url to test.")

	flag.Parse()
}

func main() {
	initConfig()
}
