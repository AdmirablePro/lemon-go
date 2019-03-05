package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var (
	logger        = logrus.New()
	taskQueue     = TaskQueue{}
	serverAddress *string
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logger.SetFormatter(customFormatter)
}

func main() {
	serverAddress = flag.String("server", "https://lemon.everyclass.xyz", "Address of lemon tree")
	localPort := flag.Int("local-port", 12345, "Port of local status server")

	logger.WithFields(logrus.Fields{"server": *serverAddress}).Info("Starting lemon")

	go fetchTask()
	go consume()

	http.HandleFunc("/", IndexHandler)
	err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(*localPort), nil)
	if err != nil {
		logger.Error(err)
	}

}
