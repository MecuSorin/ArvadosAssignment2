package main

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sending a POST request", func() {
	BeforeEach(func() {
		logDataHandled = func(int) {}
	})
	AfterEach(func() {
		logDataHandled = defaultLogDataHandled
	})
	Context("synchronously", func() {
		BeforeEach(func() {
			delayer = func(time.Duration) {}
		})

		AfterEach(func() {
			delayer = time.Sleep
		})

		Specify("Should return an error if not 200 OK", func() {
			alwaysFail, disposer := startAlwaysFailServerInstance()
			defer disposer(alwaysFail)
			packet := preparePacket(alwaysFail)
			err := sendDataToServer(packet)
			Expect(err).ShouldNot(Succeed())
		})

		Specify("Should return nil if 200 OK", func() {
			neverFail, disposer := startNeverFailServerInstance()
			defer disposer(neverFail)
			packet := preparePacket(neverFail)
			err := sendDataToServer(packet)
			Expect(err).Should(Succeed())
		})

		Specify("Should fail if it takes too long", func() {
			neverFail, disposer := startNeverFailServerInstance()
			defer disposer(neverFail)
			delayer = func(time.Duration) {
				time.Sleep(30 * time.Millisecond)
			}
			timeoutDurationForPOST = time.Millisecond
			defer func() { timeoutDurationForPOST = defaultTimeoutDurationForPOST }()
			packet := preparePacket(neverFail)
			err := sendDataToServer(packet)
			Expect(err).ShouldNot(Succeed())
		})
	})

	Context("asynchronously", func() {
		var packets []packetToDeliver

		BeforeEach(func() {
			delayer = func(time.Duration) {}

			packets = make([]packetToDeliver, 10)
			for i := 0; i < 10; i++ {
				packets[i] = packetToDeliver{dataToPOST: strings.NewReader(fmt.Sprintf("item %d", i+1))}
			}

		})

		AfterEach(func() {
			delayer = time.Sleep
			initiateContextForPacket = defaultInitiateContextForPacket
		})

		Specify("Should fail when servers allways reject", func() {
			initiateContextForPacket = func(p *packetToDeliver) func() {
				alwaysFail, disposer := startAlwaysFailServerInstance()
				p.url = alwaysFail.URL
				return func() { defer disposer(alwaysFail) }
			}
			counter := opsCounter{SendsToStop: 2}
			err := sendDataAsynch(packets, 2, &counter)
			Expect(err).ShouldNot(Succeed())
		})

		Specify("Should succeed when servers never reject", func() {
			initiateContextForPacket = func(p *packetToDeliver) func() {
				neverFail, disposer := startNeverFailServerInstance()
				p.url = neverFail.URL
				return func() { defer disposer(neverFail) }
			}
			counter := opsCounter{SendsToStop: 2}
			err := sendDataAsynch(packets[:2], 2, &counter)
			Expect(err).Should(Succeed())
		})

		Specify("Should fail when too few servers", func() {
			initiateContextForPacket = func(p *packetToDeliver) func() {
				neverFail, disposer := startNeverFailServerInstance()
				p.url = neverFail.URL
				return func() { defer disposer(neverFail) }
			}
			counter := opsCounter{SendsToStop: 5}
			err := sendDataAsynch(packets[:2], 2, &counter)
			Expect(err).ShouldNot(Succeed())
		})

		Specify("Should succeed when first few servers fails, but the last ones don't reject", func() {
			tryNumber := 0
			const failingServers = 4
			initiateContextForPacket = func(p *packetToDeliver) func() {
				tryNumber++
				if failingServers >= tryNumber {
					alwaysFail, disposer := startAlwaysFailServerInstance()
					p.url = alwaysFail.URL
					return func() { defer disposer(alwaysFail) }
				}
				neverFail, disposer := startNeverFailServerInstance()
				p.url = neverFail.URL
				return func() { defer disposer(neverFail) }
			}
			const serversToReach = 3
			counter := opsCounter{SendsToStop: serversToReach}
			err := sendDataAsynch(packets, 2, &counter)
			Expect(err).Should(Succeed())
			sends, total := counter.GetSends()
			Expect(total - sends).To(BeNumerically(">=", failingServers))
			Expect(sends).To(BeNumerically(">=", serversToReach))
		})
	})
})

func startAlwaysFailServerInstance() (*httptest.Server, func(*httptest.Server)) {
	alwaysFail := busyStatusServerHandler(0)
	server := httptest.NewServer(&alwaysFail)
	disposer := func(s *httptest.Server) {
		defer s.Close()
	}
	return server, disposer
}

func startNeverFailServerInstance() (*httptest.Server, func(*httptest.Server)) {
	neverFail := busyStatusServerHandler(100)
	server := httptest.NewServer(&neverFail)
	disposer := func(s *httptest.Server) {
		defer s.Close()
	}
	return server, disposer
}

func preparePacket(serverInstance *httptest.Server) packetToDeliver {
	return packetToDeliver{url: serverInstance.URL, dataToPOST: strings.NewReader("something")}
}
