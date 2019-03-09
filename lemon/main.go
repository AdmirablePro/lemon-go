package main

import (
	"flag"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var (
	logger    = logrus.New()
	taskQueue TaskQueue

	// below are command line parameters
	serverAddress *string

	// below are build-time variables
	ravenDSN           string
	gitRevision        string
	enableMetrics      string
	enableGlobalReport string
	defaultServer      string
	lang               string
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

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

	// init queue
	taskQueue = TaskQueue{}
	taskQueue.New()
}

func main() {
	serverAddress = flag.String("server", defaultServer, "Address of server")
	localPort := flag.Int("local-port", 12345, "Port of local status server")
	flag.Parse()

	logger.WithFields(logrus.Fields{"server": *serverAddress}).Infof(currentLangBundle.LemonStarting, gitRevision)

	go fetchTask()
	go consume()
	if enableMetrics == "true" {
		go metricsFlusher()
	}
	if enableGlobalReport == "true" {
		go globalReport()
	}

	http.HandleFunc("/", raven.RecoveryHandler(IndexHandler))
	err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(*localPort), nil)
	if err != nil {
		logger.Error(err)
	}

}
