package main

import 	"flag"

var (
	stubAcceptPercent  = flag.Int("stubPercent", 30, "Represents the chance that will accept the data")
	serverList         = flag.String("servers", "servers.json", "Provide a list of URLs where HTTP POST can be made")
	fileToUpload       = flag.String("file", "data.txt", "The data to be uploaded to the specified URLs. If no file matches will generate some random data")
	workers            = flag.Int("workers", 2, "Number of cocnurent workers that upload data")
	minimumServerCount = flag.Int("clones", 2, "The number of the servers that should receive the data file")
	isServerStub       = flag.Bool("stub", false, "If true will start as server stub and append the address in the servers file")
)

func main() {
	flag.Parse()
	
}