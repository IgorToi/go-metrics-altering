package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

const (
    pollInterval   = 2
    reportInterval = 10
	serverAdress = "http://localhost:8080"
	gaugeType = "gauge"
	countType = "count"
	PollCount = "PollCount"
)

var (
    rtm runtime.MemStats
	memory = make(map[string]float64)
	count = 0
)

func MakeRequest(url string, agent *http.Client) (string, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := agent.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
   	log.Println(string(body))
	resp.Body.Close()

	return string(body), err
}

// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>

func UrlConstructor(serverAddress string, metricType string, metricName string, metricValue interface{}) string {
	return fmt.Sprintf(fmt.Sprintf("%s/update/%s/%s/%f",serverAdress, metricType, metricName, metricValue))
}

func main() {
    // goroutine to update metrics every pollInterval
	go UpdateMetrics()   
	agent := &http.Client{}
	for {
		time.Sleep(reportInterval * time.Second)
		for i, v := range memory {
			s := UrlConstructor(serverAdress, gaugeType, i, v)   
			resp, err := MakeRequest(s, agent)
			if err != nil {
				panic(err)
			}
			fmt.Println(resp)
		}
		s := UrlConstructor(serverAdress, countType, PollCount, float64(count))   
		fmt.Println(s)
		resp, err := MakeRequest(s, agent)
		if err != nil {
			panic(err)
		}
		fmt.Println(resp)
	}   
}


func UpdateMetrics() {
	for {
        time.Sleep(pollInterval * time.Second)
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




