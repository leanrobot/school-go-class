package config

import (
	"flag"
	"fmt"
	"os"
	"time"
	log "github.com/cihub/seelog"

)

const (
	VERSION = "assignment-03.rc02"

	DEFAULT_PORT = 8080
	DEFAULT_LOG_FILE      = "etc/seelog.xml"
	DEFAULT_TEMPLATES_DIR = "src/bitbucket.org/thopet/timeserver/templates"

	DEFAULT_CHECKPOINT_INTERVAL = 10 * time.Second
	DEFAULT_AVG_RESPONSE = 700
	DEFAULT_DEVIATION = 100

	SESSION_NAME = "timeserver_css490_tompetit"
)

var (
	// the port the timeserver serves HTTP on.
	Port int

	// Flags related to communicating with the authserver.
	AuthPort int
	AuthUrl string

	// Flags related to simulating load.
	AvgResponse time.Duration
	Deviation time.Duration

	//Flags related to saving the authserver map to disk
	DumpFile string
	CheckpointInterval time.Duration

	VersionPrint bool
	TemplatesDir string
	LogConfigFile string
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
	// parse the flags and return a dictionary of all read flags.
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
	flag.DurationVar(&AvgResponse, "avg-response", DEFAULT_AVG_RESPONSE,
		`The average amount of duration in milliseconds to wait in order
		to simulate load`)
	flag.DurationVar(&Deviation, "deviation", DEFAULT_DEVIATION,
		`The value of one unit of standard deviation from the
		average response.`)

	//Flags related to saving the authserver map to disk
	flag.StringVar(&DumpFile, "dumpfile", "dumpfile.json", 
		`The location of the dumpfile for user data.`)
	//flag.DurationVar(&CheckpointInterval, )

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