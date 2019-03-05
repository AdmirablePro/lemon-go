package main

import (
	"encoding/json"
	"github.com/getsentry/raven-go"
	"net/http"
	"time"
)

// fetchTask fetches a task list from server if local size is smaller than 10 then append to local task queue.
func fetchTask() {
	for {
		if taskQueue.Size() < 10 {
			// get task
			logger.Info("Fetching task...")
			resp, err := http.Get(*serverAddress + "/task")
			if err != nil {
				raven.CaptureErrorAndWait(err, nil)
				logger.Warnf("Error when fetching task: %s", err.Error())
			}

			// decode list
			var tasks []Task
			err = json.NewDecoder(resp.Body).Decode(&tasks)
			if err != nil {
				raven.CaptureErrorAndWait(err, nil)
				logger.Warnf("Decode error when fetching task: %s", err.Error())
			}

			// save to queue
			for _, item := range tasks {
				taskQueue.Append(item)
			}
		}

		time.Sleep(time.Second) // sleep 1s
	}
}
