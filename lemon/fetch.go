package main

import (
	"encoding/json"
	"github.com/getsentry/raven-go"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// fetchTask fetches a task list from server if local size is smaller than 10 then append to local task queue.
func fetchTask(taskChannel chan<- Task) {
	for {
		time.Sleep(time.Second) // sleep 1s
		if taskList.Len() < *maxQueueSize {
			// get task
			resp, err := http.Get(*serverAddress + "/task?num=" + strconv.Itoa(*maxQueueSize-taskList.Len()))
			if err != nil {
				metricCount(fetchFailed)
				raven.CaptureErrorAndWait(err, nil)
				logger.Warnf(currentLangBundle.FetchingTaskError, err.Error())
				continue
			}

			if resp.StatusCode != 200 {
				metricCount(fetchFailed)
				raven.CaptureErrorAndWait(err, nil)
				respBody, _ := ioutil.ReadAll(resp.Body)
				logger.Warnf(currentLangBundle.FetchingTaskNon200, resp.StatusCode, string(respBody))
				continue
			}

			// decode list
			var tasks []Task
			err = json.NewDecoder(resp.Body).Decode(&tasks)
			if err != nil {
				metricCount(fetchFailed)
				raven.CaptureErrorAndWait(err, nil)
				logger.Warnf(currentLangBundle.FetchingTaskDecodeError, err.Error())
				continue
			}

			logger.Debugf(currentLangBundle.FetchTaskCount, len(tasks))

			// save to queue
			for _, item := range tasks {
				taskChannel <- item
				metricCount(taskReceived)
			}
		}
	}
}
