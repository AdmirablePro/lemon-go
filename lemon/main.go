package main

import (
	"container/list"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	logger         = logrus.New()
	taskList       *list.List
	userIdentifier string

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
	taskList = list.New()

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
	serverAddress = flag.String("server", defaultServer, "Address of server(must start with scheme)")
	localPort := flag.Int("local-port", 12345, "Port of local status server")
	flag.Parse()

	logger.WithFields(logrus.Fields{"server": *serverAddress, "user": userIdentifier}).Infof(currentLangBundle.LemonStarting, gitRevision)

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
