package main

import (
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"
)

type Sample struct {
	Time  time.Time
	Value int
}

const (
	DEF_SAMPLE_INTERVAL = 10 * time.Second
	DEF_TIMEOUT         = 5 * time.Second
	DEF_RUNTIME         = 30 * time.Second
	DEF_TARGET_STR      = "http://localhost:8080/monitor,http://localhost:9090/monitor"
)

var (
	targets        []string
	sampleInterval time.Duration
	runtime        time.Duration

	dataSet map[string]map[string][]Sample

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

func initStructure() {
	// make the whole map
	dataSet = make(map[string]map[string][]Sample)
	// make a map for every target
	for _, t := range targets {
		dataSet[t] = make(map[string][]Sample)
	}
}

func main() {
	timeout = DEF_TIMEOUT
	initFlags()
	initClient()
	initStructure()

	finish := time.Tick(runtime)
	ticker := time.Tick(sampleInterval)

monitorloop:
	for {
		// tick for collection, or exit the program. precedence is given
		// to exiting the program over ticking.
		select {
		case <-finish:
			break monitorloop
		default:
			time := <-ticker
			log.Infof("tick: %v", time)
		}

		for _, target := range targets {
			data, err := requestData(target)
			if err != nil {
				panic(err) //TODO what to do on failure to monitor?
			}
			err = recordResults(target, data)
			if err != nil {
				panic(err) //TODO what to do on failure to monitor?
			}
		}

	}

	// Print the data set collected.
	json, _ := json.MarshalIndent(dataSet, "  ", "  ")
	fmt.Println(string(json))

}

func recordResults(target string, data map[string]int) error {
	recordedTime := time.Now()
	// fmt.Printf("Target [%s]\t\t[%v]\n", target, recordedTime)
	for key, value := range data {
		if _, ok := dataSet[target][key]; !ok {
			// attempts to
			capacity := int(math.Ceil(float64(runtime / sampleInterval)))
			dataSet[target][key] = make([]Sample, 0, capacity)
		}

		sample := Sample{
			Time:  recordedTime,
			Value: value,
		}
		dataSet[target][key] = append(dataSet[target][key], sample)
	}
	return nil // TODO error handling?
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
	err = json.Unmarshal(jsonStr, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
