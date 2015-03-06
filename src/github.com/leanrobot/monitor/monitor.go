package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Sample struct {
	time  time.Time
	value int
}

const (
	DEF_SAMPLE_INTERVAL = 10 * time.Second
	DEF_TIMEOUT         = 5 * time.Second
	DEF_RUNTIME         = 30 * time.Second
	DEF_TARGET_STR      = "http://localhost:8080/monitor" //,http://localhost:9090/monitor"
)

var (
	targets        []string
	sampleInterval time.Duration
	runtime        time.Duration

	dataSet map[string][]Sample

	client  *http.Client
	timeout time.Duration
)

func initFlags() {
	targetStr := flag.String("targets", DEF_TARGET_STR, "comma separated list of urls to sample")
	targets = strings.Split(*targetStr, ",")

	flag.DurationVar(&runtime, "runtime", DEF_RUNTIME, "time to run the monitoring service")
	flag.DurationVar(&sampleInterval, "sample-interval", DEF_SAMPLE_INTERVAL,
		"the interval at which to sample")
	flag.Parse()
}

func initClient() {
	client = &http.Client{
		Timeout: timeout,
	}
}

func main() {
	dataSet = make(map[string][]Sample)
	timeout = DEF_TIMEOUT
	initFlags()
	initClient()

	ticker := time.Tick(sampleInterval)
	finish := time.Tick(runtime)

monitorloop:
	for {
		for _, target := range targets {
			data, err := requestData(target)
			if err != nil {
				panic(err) //TODO what to do on failure to monitor?
			}
			err = recordResults(data)
			if err != nil {
				panic(err) //TODO what to do on failure to monitor?
			}
		}

		// tick for collection, or exit the program. precedence is given
		// to exiting the program over ticking.
		select {
		case <-finish:
			break monitorloop
		default:
			select {
			case <-ticker:
				break
			}
		}
	}

	// print the results collected.
	keys := make([]string, len(dataSet))
	i := 0
	for key := range dataSet {
		keys[i] = key
		i += 1
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Println(key, dataSet[key])
	}
}

func recordResults(data map[string]int) error {
	for key, value := range data {
		if _, ok := dataSet[key]; !ok {
			dataSet[key] = make([]Sample, 30)
		}

		sample := Sample{
			time:  time.Now(),
			value: value,
		}
		dataSet[key] = append(dataSet[key], sample)
	}
	return nil // error handling?
}

func requestData(target string) (map[string]int, error) {
	// make http request.
	resp, err := client.Get(target)
	if err != nil {
		return nil, err
	}

	// get the body of the response as a string.
	defer resp.Body.Close()
	jsonStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// unmarshal the json
	data := make(map[string]int)
	err = json.Unmarshal(jsonStr, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
