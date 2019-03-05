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
	logger        = logrus.New()
	taskQueue     = TaskQueue{}
	serverAddress *string
	ravenDSN      string
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
	// logger.SetReportCaller(true)

	// raven
	err := raven.SetDSN(ravenDSN)
	if err != nil {
		logger.Warnf("Set DSN failed: %s", err.Error())
	}
}

func main() {
	serverAddress = flag.String("server", "https://lemon.everyclass.xyz", "Address of server")
	localPort := flag.Int("local-port", 12345, "Port of local status server")

	logger.WithFields(logrus.Fields{"server": *serverAddress}).Info("Starting lemon")

	go fetchTask()
	go consume()

	http.HandleFunc("/", raven.RecoveryHandler(IndexHandler))
	err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(*localPort), nil)
	if err != nil {
		logger.Error(err)
	}

}
