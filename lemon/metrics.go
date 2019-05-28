package main

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

const (
	FetchFailed      = "taskFetchFailedTimes" // 任务获取失败次数
	TaskReceived     = "taskReceived"         // 获取到的任务个数
	TaskSuccess      = "taskSuccess"          // 成功提交的任务个数
	TaskFailed       = "taskFailed"           // 任务失败次数
	TaskSubmitFailed = "taskSubmitFailed"     // 任务提交失败次数
)

var (
	metricMap   = map[string]*uint32{}
	metricNames = [...]string{FetchFailed, TaskReceived, TaskSuccess, TaskFailed, TaskSubmitFailed}
)

func init() {
	// init each count to zero
	for _, name := range metricNames {
		var zero = uint32(0)
		metricMap[name] = &zero
	}
}

// metricsFlusher prints metricMap and clear counts every 30 seconds by default.
func metricsFlusher(stopChan <-chan struct{}) {
	logger.Info(currentLangBundle.MetricsEnabled)

	if lightSleep(*metricsIntervalSeconds, stopChan, currentLangBundle.ExitMetricFlusher) {
		return
	}

	for {
		select {
		case <-stopChan:
			logger.Info(currentLangBundle.ExitMetricFlusher)
			return
		default:
			metricsJson, err := json.Marshal(metricMap)
			if err != nil {
				logger.Error("JSON marshal error when converting metricMap")
			}

			logger.Info(fmt.Sprintf(currentLangBundle.MetricsLogPrefix, *metricsIntervalSeconds), string(metricsJson))

			// set values to 0
			for key := range metricMap {
				atomic.StoreUint32(metricMap[key], 0)
			}

			if lightSleep(*metricsIntervalSeconds, stopChan, currentLangBundle.ExitMetricFlusher) {
				return
			}
		}
	}
}

// MetricAddOne adds 1 for the specific metric name.
func MetricAddOne(metricName string) {
	atomic.AddUint32(metricMap[metricName], 1)
}

// globalReport gets global status of current project from server and print it to console.
// TODO: implement this
func globalReport() {
	logger.Info(currentLangBundle.GlobalReportEnabled)

	// get status from server

	// json unmarshal

	// print in log
}
