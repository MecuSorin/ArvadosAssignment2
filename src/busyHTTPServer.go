package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// busyStatusServerHandler represents the percent chance that the POST call will succeed
type busyStatusServerHandler int

var (
	delayer = time.Sleep


	delayerSleepDuration = time.Duration(5 * time.Second)
)

// serveHTTP waits then will return a random success/fail header according to
// the value of responseStatus
func (responseStatus *busyStatusServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	delayer(delayerSleepDuration)

	returnHeaderStatus := http.StatusNotFound
	if responseStatus.getRandomChanceOfSuccess() {
		returnHeaderStatus = http.StatusOK
	}

	w.WriteHeader(returnHeaderStatus)
	logDataHandled(returnHeaderStatus)
}

// getRandomChanceOfSuccess provides a random boolean based on responseStatus percent chance
func (responseStatus *busyStatusServerHandler) getRandomChanceOfSuccess() bool {
	return rand.Intn(100) < int(*responseStatus)
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
