package main

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// consume gets task from local queue and do the task.
func consume() {
	for {
		task := taskQueue.Pop()
		var (
			resp *http.Response
			err  error
		)
		if task.Method == "POST" {
			// resp, err = http.Post(task.Host + task.Path)
		} else if task.Method == "GET" {
			resp, err = http.Get(task.Host + task.Path)
		} else {
			logger.WithFields(logrus.Fields{"taskID": task.TaskID}).Warn("HTTP method not supported. Ignore task.")
			continue
		}
		if err != nil {
			logger.Warnf("Error when consuming task: %s", err.Error())
			continue
		}

		// read body
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Warnf("Error when reading body: %s", err.Error())
			continue
		}

		// make result
		result := Result{
			Status:       "",
			TaskID:       task.TaskID,
			ResponseCode: resp.StatusCode,
			Data:         string(bodyBytes),
			FetchedTime:  "",
			UserAgent:    "Go client"}
		resultBytes, err := json.Marshal(result)
		if err != nil {
			logger.Warnf("Error when marshaling result: %s", err.Error())
			continue
		}

		// post result to server
		resp, err = http.Post("", "application/json;charset=utf-8", bytes.NewBuffer(resultBytes))
		if err != nil {
			logger.Warnf("Error when posting result to server: %s", err.Error())
			continue
		}
	}
}
