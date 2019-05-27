package main

import (
	"time"
)

// lightSleep checks channel every second while sleeping. Returns true if be waken up.
func lightSleep(sleepSeconds int, stopChan <-chan struct{}, exitLog string) bool {
	var i = 0
	for i < sleepSeconds {
		select {
		case <-stopChan:
			logger.Info(exitLog)
			return true
		default:
			time.Sleep(time.Second * time.Duration(1))
			i += 1
		}
	}
	return false
}
