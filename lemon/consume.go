package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

// consume gets task from local queue and do the task.
func consume() {
	for {
		task := taskQueue.Pop()
		if task == nil {
			continue
		}
		var (
			resp *http.Response
			err  error
		)
		if task.HTTPMethod == "POST" {
			// resp, err = http.Post(task.Host + task.Path)
			// todo: implement post
		} else if task.HTTPMethod == "GET" {
			resp, err = http.Get(task.Host + task.Path)
			// todo: add header, cookie, etc.
			// todo: 指数等待
			// todo: 怎么定义成功失败？
		} else {
			metricCount(M_TASK_FAILED)
			logger.WithFields(logrus.Fields{"taskID": task.TaskID}).Warn("HTTP method not supported. Ignore task.")
			continue
		}
		if err != nil {
			metricCount(M_TASK_FAILED)
			raven.CaptureErrorAndWait(err, nil)
			logger.Warnf("Error when consuming task: %s", err.Error())
			continue
		}

		// read body
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			metricCount(M_TASK_FAILED)
			raven.CaptureErrorAndWait(err, nil)
			logger.Warnf("Error when reading body: %s", err.Error())
			continue
		}

		// check task status
		var taskStatus string
		if resp.StatusCode == 200 {
			taskStatus = "success"
		} else {
			taskStatus = "error"
		}

		// make result
		result := Result{
			Status:       taskStatus,
			TaskID:       task.TaskID,
			ResponseCode: resp.StatusCode,
			Data:         string(bodyBytes),
			FetchedTime:  time.Now().Unix(),
			UserAgent:    fmt.Sprintf("Go client(%s)", gitRevision)}
		resultBytes, err := json.Marshal(result)
		if err != nil {
			metricCount(M_TASK_FAILED)
			raven.CaptureErrorAndWait(err, nil)
			logger.Warnf("Error when marshaling result: %s", err.Error())
			continue
		}

		// post result to server
		resp, err = http.Post("", "application/json;charset=utf-8", bytes.NewBuffer(resultBytes))
		if err != nil {
			metricCount(M_TASK_FAILED)
			raven.CaptureErrorAndWait(err, nil)
			logger.Warnf("Error when posting result to server: %s", err.Error())
			continue
		}
		metricCount(M_TASK_SUCCESS)
	}
}
