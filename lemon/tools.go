package main

import (
	"time"
)

// lightSleep checks channel every second while sleeping. Returns true if be waken up.
func lightSleep(sleepSeconds int, stopChan <-chan struct{}) bool {
	var i = 0
	for i < sleepSeconds {
		select {
		case <-stopChan:
			logger.Info("Exit metrics flusher")
			return true
		default:
			time.Sleep(time.Second * time.Duration(1))
			i += 1
		}
	}
	return false
}
