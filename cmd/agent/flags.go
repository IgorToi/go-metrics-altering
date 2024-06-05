package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)


var (
	flagRunAddr string
	flagReportInterval int
	flagPollInterval int
	err error
)

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", 10, "frequency of metrics being sent to the server")
	flag.IntVar(&flagPollInterval, "p", 2, "frequency of metrics being received from the runtime package")

	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
        flagRunAddr = envRunAddr
	}
	if envRoportInterval := os.Getenv("REPORT_INTERVAL"); envRoportInterval != "" {
        flagReportInterval, err = strconv.Atoi(envRoportInterval) 
		if err != nil {
			log.Fatal(err)
		}
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
        flagPollInterval, err = strconv.Atoi(envPollInterval) 
		if err != nil {
			log.Fatal(err)
		}
	}
}