package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

const (
	M_FETCH_FAILED  = "taskFetchFailedTimes" // 任务获取失败次数
	M_TASK_RECEIVED = "taskReceived"         // 获取到的任务个数
	M_TASK_SUCCESS  = "taskSuccess"          // 成功提交的任务个数
	M_TASK_FAILED   = "taskFailed"           // 任务失败次数
)

var (
	metrics      = map[string]int{M_FETCH_FAILED: 0, M_TASK_RECEIVED: 0, M_TASK_SUCCESS: 0, M_TASK_FAILED: 0}
	metricsMutex = sync.RWMutex{}
)

// metricsFlusher prints metrics every 30 seconds and clear counts.
func metricsFlusher() {
	interval := 30
	time.Sleep(time.Second * time.Duration(interval))
	for {
		metricsJson, err := json.Marshal(metrics)
		if err != nil {
			logger.Error("JSON marshal error when converting metrics")
		}

		logger.Info(fmt.Sprintf("Metrics in last %ds: ", interval), string(metricsJson))
		metricsMutex.Lock()
		for key := range metrics {
			metrics[key] = 0
		}
		metricsMutex.Unlock()
		time.Sleep(time.Second * time.Duration(interval))
	}
}

// add 1 for the specific metric name
func metricCount(metricName string) {
	metricsMutex.Lock()
	metrics[metricName] += 1
	metricsMutex.Unlock()
}
