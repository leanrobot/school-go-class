package config

import (
	"flag"
	"fmt"
	log "github.com/cihub/seelog"
	"os"
)

const (
	VERSION = "assignment-04.rc01"

	DEFAULT_PORT          = 8080
	DEFAULT_LOG_FILE      = "etc/timeserver_seelog.xml"
	DEFAULT_TEMPLATES_DIR = "src/github.com/leanrobot/timeserver/templates"

	DEFAULT_CHECKPOINT_INTERVAL = 10
	DEFAULT_AVG_RESPONSE        = 5000
	DEFAULT_DEVIATION           = 500
	DEFAULT_AUTH_TIMEOUT        = 1000

	SESSION_NAME = "timeserver_css490_tompetit"

	DEFAULT_MAX_REQUESTS = 0
)

var (
	// the port the timeserver serves HTTP on.
	Port int

	// Flags related to communicating with the authserver.
	AuthPort int
	AuthUrl  string

	// Flags related to simulating load.
	AvgResponse  int
	Deviation    int
	RequestLimit int

	//Flags related to saving the authserver map to disk
	DumpFile           string
	CheckpointInterval int

	VersionPrint  bool
	TemplatesDir  string
	LogConfigFile string

	AuthTimeout int
)

func init() {
	initFlags()
	initLogger(LogConfigFile)
	if VersionPrint {
		fmt.Println(VERSION)
		os.Exit(0)
	}
}

func initFlags() {
	flag.IntVar(&Port, "port", DEFAULT_PORT,
		"port to launch webserver on, default is 8080")
	flag.BoolVar(&VersionPrint, "V", false,
		"Display version information")
	flag.StringVar(&TemplatesDir, "templates", DEFAULT_TEMPLATES_DIR,
		"the location of site templates")
	flag.StringVar(&LogConfigFile, "log", DEFAULT_LOG_FILE,
		"the location of the seelog configuration file")
	flag.StringVar(&AuthUrl, "authhost", "localhost",
		"The network address for the auth server")
	flag.IntVar(&AuthPort, "authport", 9090,
		"The port which to connect to the authserver on.")

	// Flags related to simulating load.
	flag.IntVar(&AvgResponse, "avg-response-ms", DEFAULT_AVG_RESPONSE,
		`The average amount of duration in milliseconds to wait in order
		to simulate load`)
	flag.IntVar(&Deviation, "deviation-ms", DEFAULT_DEVIATION,
		`The value of one unit of standard deviation from the
		average response.`)

	// Timeout
	flag.IntVar(&AuthTimeout, "auth-timeout-ms", DEFAULT_AUTH_TIMEOUT,
		"The timeout in milliseconds when timeserver talks to the authserver.")

	//Flags related to saving the authserver map to disk
	flag.StringVar(&DumpFile, "dumpfile", "",
		`The location of the dumpfile for user data.`)
	flag.IntVar(&CheckpointInterval, "checkpoint-interval-ms",
		DEFAULT_CHECKPOINT_INTERVAL,
		"Performs a save to dumpfile every checkpoint-interval.")

	//Flags for request limiting
	flag.IntVar(&RequestLimit, "max-inflight", 0,
		"The maximum amount of conurrent requests to serve.")

	flag.Parse()
}

func initLogger(configFile string) {
	// Setup the logger as the default package logger.
	logger, err := log.LoggerFromConfigAsFile(configFile)
	if err != nil {
		panic(err)
	}
	log.ReplaceLogger(logger)
	defer log.Flush()
}
