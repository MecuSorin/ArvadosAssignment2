package main

import (
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("busyHTTPServer stub ", func() {

	Specify("Should wait a random time on each request/response", func() {
		delayerDone := make(chan bool)
		delayer = func(time.Duration) {
			delayerDone <- true
		}
		defer func() { delayer = time.Sleep }()

		randomChanceToSucceed := busyStatusServerHandler(50)
		serverInstance := httptest.NewServer(&randomChanceToSucceed)
		defer serverInstance.Close()
		go func() {
			<-delayerDone
		}()
		packet := preparePacket(serverInstance)
		sendDataToServer(packet)
	})
})
