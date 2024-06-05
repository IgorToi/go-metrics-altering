package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
    pollInterval   = 2
	gaugeType = "gauge"
	countType = "counter"
	PollCount = "PollCount"
	StatusOK = 200
	ProtocolScheme = "http://"
)

var (
    rtm runtime.MemStats
	memory = make(map[string]float64)
	count = 0
)

func SendMetric(requestURL, metricType, metricName, metricValue string, req *resty.Request) (*resty.Response, error) {
	return req.Post(req.URL + "/update/{metricType}/{metricName}/{metricValue}")
}

// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>

func main() {
	parseFlags()
    // start goroutine to update metrics every pollInterval
	go UpdateMetrics()   
	agent := resty.New()
	for {
		time.Sleep(time.Duration(flagReportInterval) * time.Second)
		for i, v := range memory {
				req := agent.R()
				req.SetPathParams(map[string]string{
					"metricType": gaugeType,
					"metricName": i,
					"metricValue": strconv.FormatFloat(v, 'f', 6, 64 ),
				}).SetHeader("Content-Type", "text/plain")
				req.URL = ProtocolScheme + flagRunAddr
				_, err := SendMetric(req.URL, gaugeType, i,  strconv.FormatFloat(v, 'f', 6, 64 ), req)
				if err != nil {
					panic(err)
				}
				fmt.Println("Metric has been sent successfully")
		}
		req := agent.R()
		req.SetPathParams(map[string]string{
			"metricType": countType,
			"metricName": PollCount,
			"metricValue": strconv.Itoa(count),
		}).SetHeader("Content-Type", "text/plain")
		req.URL = ProtocolScheme + flagRunAddr
		_, err := SendMetric(req.URL, countType, PollCount,  strconv.Itoa(count), req)
		if err != nil {
			panic(err)
		}
		fmt.Println("Metric has been sent successfully")
	}   
}

func UpdateMetrics() {
	for {
        time.Sleep(time.Duration(flagPollInterval) * time.Second)
        runtime.ReadMemStats(&rtm)

		memory["Alloc"] = float64(rtm.Alloc)
		memory["BuckHashSys"] = float64(rtm.BuckHashSys)
		memory["Frees"] = float64(rtm.Frees)
		memory["GCCPUFraction"] = float64(rtm.GCCPUFraction )
		memory["GCSys"] = float64(rtm.GCSys)
		memory["HeapAlloc"] = float64(rtm.HeapAlloc)
		memory["HeapIdle"] = float64(rtm.HeapIdle)
		memory["HeapInuse"] = float64(rtm.HeapInuse)
		memory["HeapObjects"] = float64(rtm.HeapObjects)
		memory["HeapReleased"] = float64(rtm.HeapReleased)
		memory["HeapSys"] = float64(rtm.HeapSys)
		memory["LastGC"] = float64(rtm.LastGC)
		memory["Lookups"] = float64(rtm.Lookups)
		memory["MCacheInuse"] = float64(rtm.MCacheInuse)
		memory["MCacheSys"] = float64(rtm.MCacheSys)
		memory["MSpanInuse"] = float64(rtm.MSpanInuse)
		memory["MSpanSys"] = float64(rtm.MSpanSys)
		memory["Mallocs"] = float64(rtm.Mallocs)
		memory["NextGC"] = float64(rtm.NextGC)
		memory["NumForcedGC"] = float64(rtm.NumForcedGC)
		memory["NumGC"] = float64(rtm.NumGC)
		memory["OtherSys"] = float64(rtm.OtherSys)
		memory["NextGC"] = float64(rtm.NextGC)
		memory["NumForcedGC"] = float64(rtm.NumForcedGC)
		memory["NumGC"] = float64(rtm.NumGC)
		memory["OtherSys"] = float64(rtm.OtherSys)
		memory["PauseTotalNs"] = float64(rtm.PauseTotalNs)
		memory["StackInuse"] = float64(rtm.StackInuse)
		memory["StackSys"] = float64(rtm.StackSys)
		memory["Sys"] = float64(rtm.StackSys)
		memory["TotalAlloc"] = float64(rtm.TotalAlloc)
		memory["RandomValue"] = rand.Float64()
		count++
    }
}




