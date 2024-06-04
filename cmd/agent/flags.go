package main

import "flag"


var (
	flagRunAddr string
	flagReportInterval int
	flagPollInterval int
)

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", 10, "frequency of metrics being sent to the server")
	flag.IntVar(&flagPollInterval, "p", 2, "frequency of metrics being received from the runtime package")

	flag.Parse()
}