package main

import (
	"encoding/json"
	"flag"
	"github.com/getsentry/raven-go"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
)

var (
	logger = logrus.New()

	userIdentifier string

	// command line parameters
	serverAddress          *string
	maxQueueSize           *int
	metricsIntervalSeconds *int

	// build-time variables
	ravenDSN           string
	gitRevision        string
	enableMetrics      string
	enableGlobalReport string
	defaultServer      string
	lang               string
)

func init() {
	// logrus formatter
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logger.SetFormatter(customFormatter)
	logger.SetLevel(logrus.DebugLevel)

	// raven
	err := raven.SetDSN(ravenDSN)
	if err != nil {
		logger.Warnf("Set DSN failed: %s", err.Error())
	}

	// user identifier
	userIdentifier = getUserIdentifier()
}

// getUserIdentifier reads lemon_seed from current directory, if no file exists, generate one.
func getUserIdentifier() string {
	var config Configuration
	if _, err := os.Stat("./lemon_seed"); os.IsNotExist(err) {
		// not exist, generate one
		config.ClientID = uuid.NewV4().String()
		configBytes, err := json.Marshal(config)
		if err != nil {
			logger.Error("Marshal error when generating configBytes: %s", err.Error())
		}

		err = ioutil.WriteFile("./lemon_seed", configBytes, 0666)
		if err != nil {
			logger.Error("Error when write config to seed: %s", err.Error())
		}
		return config.ClientID
	} else {
		// lemon_config exist
		body, err := ioutil.ReadFile("./lemon_seed")
		if err != nil {
			logger.Fatal("Read exception.")
		}
		err = json.Unmarshal(body, &config)
		if err != nil {
			logger.Fatalf("Unmarshal error when reading seed: %s", err.Error())
		}
		return config.ClientID
	}
}

func main() {
	var stopChannel = make(chan struct{})
	var wg sync.WaitGroup

	serverAddress = flag.String("server", defaultServer, "Address of server(must start with scheme)")
	maxQueueSize = flag.Int("queue-size", 10, "Max queue size")
	metricsIntervalSeconds = flag.Int("metrics-interval", 30, "Metrics interval")
	flag.Parse()

	logger.WithFields(logrus.Fields{
		"server":    *serverAddress,
		"user":      userIdentifier,
		"queueSize": *maxQueueSize}).Infof(currentLangBundle.Starting, gitRevision)

	taskChannel := make(chan Task)

	// task fetching goroutine
	go func(stop <-chan struct{}) {
		defer wg.Done()
		wg.Add(1)
		fetchTask(taskChannel, stop)
	}(stopChannel)

	// task consuming goroutine
	go func(stop <-chan struct{}) {
		defer wg.Done()
		wg.Add(1)
		consume(taskChannel, stop)
	}(stopChannel)

	// metrics goroutine
	if enableMetrics == "true" {
		go func(stop <-chan struct{}) {
			defer wg.Done()
			wg.Add(1)
			metricsFlusher(stop)
		}(stopChannel)
	}

	// global report goroutine
	if enableGlobalReport == "true" {
		go globalReport()
	}

	// handle signal
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)
	<-signalChannel // block until receive quit signal

	logger.Info(currentLangBundle.Exiting)

	// notify each goroutine to exit
	close(stopChannel)

	// wait until all goroutine exit
	wg.Wait()
	logger.Info(currentLangBundle.Exited)
}
