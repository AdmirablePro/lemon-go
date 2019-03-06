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
	client := &http.Client{}
	sleepSeconds := 1

	for {
		// sleep
		time.Sleep(time.Second * time.Duration(sleepSeconds))

		task := taskQueue.Pop()
		if task == nil {
			// if currently no tasks, sleep for 1 second.
			time.Sleep(time.Second)
			continue
		}
		var (
			resp    *http.Response
			request *http.Request
			err     error
		)

		// handle unsupported HTTP Methods
		if task.HTTPMethod != "POST" && task.HTTPMethod != "GET" {
			metricCount(M_TASK_FAILED)
			logger.WithFields(logrus.Fields{"taskID": task.TaskID}).Warn("HTTP method not supported. Ignore task.")
			continue
		}

		request, err = http.NewRequest(task.HTTPMethod, task.Host+task.Path, bytes.NewBuffer([]byte(task.Payload)))
		if err != nil {
			metricCount(M_TASK_FAILED)
			raven.CaptureErrorAndWait(err, nil)
			logger.Warnf("Error when building request: %s", err.Error())
			continue
		}

		// add header and cookie
		for k, v := range task.Header {
			request.Header.Set(k, v)
		}
		request.Header.Set("Cookie", task.Cookie)

		// build query string
		q := request.URL.Query()
		for k, v := range task.Param {
			q.Set(k, v)
		}
		request.URL.RawQuery = q.Encode()

		// do request
		resp, err = client.Do(request)
		if err != nil {
			sleepSeconds *= 2 // add time of sleep when request fails
			metricCount(M_TASK_FAILED)
			raven.CaptureErrorAndWait(err, nil)
			logger.Warnf("Error when consuming task: %s", err.Error())
			continue
		}

		// read body
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close() //must close resp.Body
		if err != nil {
			metricCount(M_TASK_FAILED)
			raven.CaptureErrorAndWait(err, nil)
			logger.Warnf("Error when reading response body: %s", err.Error())
			continue
		}

		// check task status
		// todo: 怎么定义成功失败？
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

		// reduce time of sleep after success
		if sleepSeconds > 1 {
			sleepSeconds /= 2
		}
	}
}
