package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	defaultInitiateContextForPacket = func(*packetToDeliver) func() { return func() {} }
	defaultTimeoutDurationForPOST   = time.Duration(6 * time.Second)

	initiateContextForPacket = defaultInitiateContextForPacket
	timeoutDurationForPOST   = defaultTimeoutDurationForPOST
)

type packetToDeliver struct {
	dataToPOST io.Reader
	url        string
	filePath   string
}

// sendDataToServer sends a packet to a server and returns nil if succeeded otherwise an error
func sendDataToServer(packet packetToDeliver) error {
	postChannel := make(chan error, 1)
	go func() {
		resp, err := http.Post(packet.url, "application/octet-stream", packet.dataToPOST)
		if nil != err {
			postChannel <- err
			return
		}
		defer resp.Body.Close()
		if http.StatusOK != resp.StatusCode {
			postChannel <- fmt.Errorf("Failed to send the data: %q", resp.Status)
			return
		}
		postChannel <- nil
	}()
	select {
	case result := <-postChannel:
		return result
	case <-time.After(timeoutDurationForPOST):
		return fmt.Errorf("Failed to send the data: Exceeded the timeout period by %v ", timeoutDurationForPOST)
	}
}

// sender is blocking until a job is available, then is making the POST requests and
// blocks until can process the result
func sender(packets <-chan packetToDeliver, counter *opsCounter, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	for senderJob := range packets {
		disposeContext := initiateContextForPacket(&senderJob)
		result := sendDataToServer(senderJob)
		disposeContext()
		counter.Process(result)
		if counter.Done() {
			break
		}
	}
}

// sendDataAsynch tries the best to send the data to atleast counter.SendsToStop servers
func sendDataAsynch(packets []packetToDeliver, workersCount int, counter *opsCounter) error {
	const MaxJobs = 100
	jobs := make(chan packetToDeliver, MaxJobs)
	var waitGroup sync.WaitGroup
	// start workers that blocks until job available
	for worker := 0; worker < workersCount; worker++ {
		waitGroup.Add(1)
		go sender(jobs, counter, &waitGroup)
	}

	// until done: lazy feed the workers with jobs (blocks when no worker available)
	go func() {
		defer close(jobs)

		for _, packet := range packets {
			jobs <- packet
		}
	}()
	waitGroup.Wait()
	if !counter.Done() {
		sends, total := counter.GetSends()
		return fmt.Errorf("Failed: Managed to send just to %d/%d servers", sends, total)
	}
	return nil
}
