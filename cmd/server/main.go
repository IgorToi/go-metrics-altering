package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauge     map[string]float64
	counter   map[string]int64
}

var memory = &MemStorage{
	gauge: make(map[string]float64),
	counter: make(map[string]int64),
}

func reqeustHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	pathSlice := strings.Split(r.URL.Path, "/")
	if len(pathSlice) <= 3 {
		fmt.Println(pathSlice)
		w.WriteHeader(http.StatusNotFound)
        return
	} 
	if pathSlice[2] != "gauge" && pathSlice[2] != "counter" {
		w.WriteHeader(http.StatusExpectationFailed)
        return
	}

	switch pathSlice[2] {
	case "gauge":
		if value, err := strconv.ParseFloat(pathSlice[4], 64); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			memory.gauge[pathSlice[2]] = value
		}
	case "counter": 
		if value, err := strconv.ParseInt(pathSlice[4], 0, 64); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			memory.counter[pathSlice[3]] += value
		}
	}
}

func main() {
	mux := http.NewServeMux()

	// endpoint `/update`
	mux.HandleFunc(`/update/`, reqeustHandler) 

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
