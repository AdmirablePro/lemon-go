package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

var logger = logrus.New()

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
	logger.WithFields(logrus.Fields{"port": 12345}).Info("Starting lemon")
	http.HandleFunc("/", IndexHandler)
	err := http.ListenAndServe("127.0.0.1:12345", nil)
	if err != nil {
		logger.Error(err)
	}

}
