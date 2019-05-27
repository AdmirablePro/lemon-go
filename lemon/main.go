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
	var stopChannels []chan struct{}
	var wg sync.WaitGroup

	serverAddress = flag.String("server", defaultServer, "Address of server(must start with scheme)")
	maxQueueSize = flag.Int("queue-size", 10, "Max queue size")
	metricsIntervalSeconds = flag.Int("metrics-interval", 30, "Metrics interval")
	flag.Parse()

	logger.WithFields(logrus.Fields{
		"server":    *serverAddress,
		"user":      userIdentifier,
		"queueSize": *maxQueueSize}).Infof(currentLangBundle.LemonStarting, gitRevision)

	taskChannel := make(chan Task)

	stopChan := make(chan struct{})
	stopChannels = append(stopChannels, stopChan)
	go func(stop <-chan struct{}) {
		defer wg.Done()
		wg.Add(1)
		fetchTask(taskChannel, stop)
	}(stopChan)

	go consume(taskChannel)

	if enableMetrics == "true" {
		stopChan := make(chan struct{})
		stopChannels = append(stopChannels, stopChan)
		go func(stop <-chan struct{}) {
			defer wg.Done()
			wg.Add(1)
			metricsFlusher(stop)
		}(stopChan)
	}
	if enableGlobalReport == "true" {
		go globalReport()
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)
	<-signalChannel // block until receive quit signal
	//todo: exit each goroutines
	logger.Info(currentLangBundle.Exiting)

	// notify each goroutine to exit
	for _, ch := range stopChannels {
		close(ch)
	}

	// wait until all goroutine exit
	wg.Wait()
	logger.Info(currentLangBundle.Exited)
}
