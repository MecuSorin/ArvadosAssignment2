package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
)

var (
	stubAcceptPercent  = flag.Int("stubPercent", 30, "Represents the chance that will accept the data")
	serverList         = flag.String("servers", "servers.json", "Provide a list of URLs where HTTP POST can be made")
	fileToUpload       = flag.String("file", "data.txt", "The data to be uploaded to the specified URLs. If no file matches will generate some random data")
	workers            = flag.Int("workers", 2, "Number of coccurent workers that upload data")
	minimumServerCount = flag.Int("clones", 2, "The number of the servers that should receive the data file")
	isServerStub       = flag.Bool("stub", false, "If true will start as server stub and append the address in the servers file")
)

func main() {
	flag.Parse()

	if *isServerStub {
		stubPercent := min(100, max(0, *stubAcceptPercent))
		serverHandler := busyStatusServerHandler(stubPercent)
		s := httptest.NewServer(&serverHandler)
		defer s.Close()
		urls, err := readServers(*serverList)
		if nil != err {
			urls = []string{}
		}
		if !contains(urls, s.URL) {
			urls = append(urls, s.URL)
			writeServers(*serverList, urls)
		}

		fmt.Printf("Listening on %q. Having %d %% accept. Press any key ...", s.URL, stubPercent)
		r := bufio.NewReader(os.Stdin)
		r.ReadByte()
		return
	}

	urls, err := readServers(*serverList)
	crashOnError(err)
	packets := make([]packetToDeliver, len(urls))
	for url := range urls {
		packets[url].url = urls[url]
		packets[url].filePath = *fileToUpload
	}
	wokersNo := min(max(1, *workers), 100)
	serversToUse := min(10, max(1, *minimumServerCount))
	fmt.Printf("Sending data using %d workers to %d servers\n", wokersNo, serversToUse)
	err = sendDataAsynch(packets, wokersNo, &opsCounter{SendsToStop: serversToUse})
	if nil == err {
		fmt.Println("Done! Press any key ...")
		r := bufio.NewReader(os.Stdin)
		r.ReadByte()
		return
	}
	fmt.Println(err.Error())
	fmt.Println("Press any key ...")
	r := bufio.NewReader(os.Stdin)
	r.ReadByte()
	return
}

func readServers(filePath string) ([]string, error) {
	var data []string
	file, err := ioutil.ReadFile(filePath)
	if nil != err {
		return nil, err
	}
	crashOnError(json.Unmarshal(file, &data))
	return data, nil
}

func writeServers(filePath string, servers []string) {
	file, err := json.Marshal(servers)
	crashOnError(err)
	crashOnError(ioutil.WriteFile(filePath, file, 0644))
}

func crashOnError(err error) {
	if nil != err {
		log.Fatal(err)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
