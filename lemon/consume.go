package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	taskStatusError              = "error"
	taskStatusSuccess            = "success"
	errorCodeUnsupportedMethod   = 4001
	errorCodeUnsupportedScheme   = 4002
	errorCodeFailBuildingRequest = 4003
	errorCodeRequestFailed       = 5001
	errorCodeReadBodyFailed      = 5002
)

// makeErrorResult makes a Result object for error.
func makeErrorResult(taskID string, errorCode int) *Result {
	result := Result{
		Status:       taskStatusError,
		TaskID:       taskID,
		FetchedTime:  time.Now().Unix(),
		ResponseCode: 0,
		Data:         "",
		User:         fmt.Sprintf("Go client(%s)", userIdentifier),
		ErrorCode:    errorCode}
	return &result
}

// post a `Result` to lemon server
func report(result *Result) {
	// marshal
	resultBytes, err := json.Marshal(result)
	if err != nil {
		MetricAddOne(TaskSubmitFailed)
		raven.CaptureErrorAndWait(err, nil)
		logger.Warnf("Error when marshaling result: %s", err.Error())
		return
	}

	// post result to server
	resp, err := http.Post(*serverAddress+"/task", "application/json;charset=utf-8", bytes.NewBuffer(resultBytes))
	if err != nil {
		MetricAddOne(TaskSubmitFailed)
		raven.CaptureErrorAndWait(err, nil)
		logger.Warnf(currentLangBundle.SubmitResultError, err.Error())
		return
	}

	if resp.StatusCode != 200 {
		MetricAddOne(TaskSubmitFailed)
		respBody, _ := ioutil.ReadAll(resp.Body)
		logger.Warnf(currentLangBundle.SubmitResultNon200, resp.StatusCode, respBody)
	} else {
		MetricAddOne(TaskSuccess)
	}
}

// make a `Result` for error and report to server
func reportError(taskID string, errorCode int) {
	result := makeErrorResult(taskID, errorCode)
	report(result)
}

// consume gets task from local queue and do the task.
func consume(taskChannel <-chan Task) {
	client := &http.Client{}
	sleepSeconds := 1

	for {
		// sleep between each requests
		time.Sleep(time.Second * time.Duration(sleepSeconds))

		task := <-taskChannel
		var (
			resp    *http.Response
			request *http.Request
			err     error
		)

		// handle unsupported HTTP Methods
		if task.HTTPMethod != "POST" && task.HTTPMethod != "GET" {
			MetricAddOne(TaskFailed)
			logger.WithFields(logrus.Fields{"taskID": task.TaskID}).Warn("HTTP method not supported. Ignore task.")
			reportError(task.TaskID, errorCodeUnsupportedMethod)
			continue
		}

		if scheme := strings.ToLower(task.Scheme); scheme != "http" && scheme != "https" {
			MetricAddOne(TaskFailed)
			logger.WithFields(logrus.Fields{"taskID": task.TaskID}).Warnf("Scheme %s not supported. Ignore task.", task.Scheme)
			reportError(task.TaskID, errorCodeUnsupportedScheme)
			continue
		}

		request, err = http.NewRequest(task.HTTPMethod, task.Scheme+"://"+task.Host+task.Path, bytes.NewBuffer([]byte(task.Payload)))
		if err != nil {
			MetricAddOne(TaskFailed)
			raven.CaptureErrorAndWait(err, nil)
			logger.Warnf("Error when building request: %s", err.Error())
			reportError(task.TaskID, errorCodeFailBuildingRequest)
			continue
		}

		// add header and cookie
		for k, v := range task.Header {
			request.Header.Set(k, v)
		}
		if task.Cookie != "" {
			request.Header.Set("Cookie", task.Cookie)
		}

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
			MetricAddOne(TaskFailed)
			raven.CaptureErrorAndWait(err, nil)
			logger.Warnf(currentLangBundle.ConsumingHTTPDoError, err.Error())
			reportError(task.TaskID, errorCodeRequestFailed)
			continue
		}

		// read body
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close() //must close resp.Body
		if err != nil {
			MetricAddOne(TaskFailed)
			raven.CaptureErrorAndWait(err, nil)
			logger.Warnf("Error when reading response body: %s", err.Error())
			reportError(task.TaskID, errorCodeReadBodyFailed)
			continue
		}

		// make result
		result := Result{
			Status:       taskStatusSuccess,
			TaskID:       task.TaskID,
			CookieID:     task.CookieID,
			ResponseCode: resp.StatusCode,
			Data:         string(bodyBytes),
			FetchedTime:  time.Now().Unix(),
			User:         fmt.Sprintf("Go client(%s)", userIdentifier)}
		report(&result)

		// reduce time of sleep after success
		if sleepSeconds > 1 {
			sleepSeconds /= 2
		}
	}
}
