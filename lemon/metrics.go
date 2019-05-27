package main

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"
)

const (
	fetchFailed      = "taskFetchFailedTimes" // 任务获取失败次数
	taskReceived     = "taskReceived"         // 获取到的任务个数
	taskSuccess      = "taskSuccess"          // 成功提交的任务个数
	taskFailed       = "taskFailed"           // 任务失败次数
	taskSubmitFailed = "taskSubmitFailed"     // 任务提交失败次数
)

var (
	metricMap   = map[string]*uint32{}
	metricNames = [...]string{fetchFailed, taskReceived, taskSuccess, taskFailed, taskSubmitFailed}
)

func init() {
	// init each count to zero
	for _, name := range metricNames {
		var zero = uint32(0)
		metricMap[name] = &zero
	}
}

// metricsFlusher prints metricMap and clear counts every 30 seconds by default.
func metricsFlusher() {
	logger.Info(currentLangBundle.MetricsEnabled)

	time.Sleep(time.Second * time.Duration(*metricsIntervalSeconds))
	for {
		metricsJson, err := json.Marshal(metricMap)
		if err != nil {
			logger.Error("JSON marshal error when converting metricMap")
		}

		logger.Info(fmt.Sprintf(currentLangBundle.MetricsInLog, *metricsIntervalSeconds), string(metricsJson))

		// set values to 0
		for key := range metricMap {
			atomic.StoreUint32(metricMap[key], 0)
		}
		time.Sleep(time.Second * time.Duration(*metricsIntervalSeconds))
	}
}

// metricCount adds 1 for the specific metric name.
func metricCount(metricName string) {
	atomic.AddUint32(metricMap[metricName], 1)
}

// globalReport uploads local statistics to central server.
// TODO: implement this
func globalReport() {
	logger.Info(currentLangBundle.GlobalReportEnabled)

	// get status from server

	// json unmarshal

	// print in log
}
