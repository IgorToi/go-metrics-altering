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
	pathSlice := pathCleaner(r.URL.Path) 
	
	
	fmt.Println(pathSlice)

	if len(pathSlice) <= 3 {
		fmt.Println(pathSlice)
		w.WriteHeader(http.StatusNotFound)
        return
	} 
	if pathSlice[1] != "gauge" && pathSlice[1] != "counter" {
		w.WriteHeader(http.StatusBadRequest)
        return
	}

	switch pathSlice[1] {
	case "gauge":
		if value, err := strconv.ParseFloat(pathSlice[3], 64); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			memory.gauge[pathSlice[2]] = value
		}
	case "counter": 
		if value, err := strconv.ParseInt(pathSlice[3], 0, 64); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			memory.counter[pathSlice[3]] += value
		}
	}
}

func pathCleaner(path string) []string {
	path = strings.TrimSpace(path)
	if strings.HasPrefix(path, "/") {
        path = path[1:]
    }
	if strings.HasSuffix(path, "/") {
        cut_off_last_char_len := len(path) - 1
        path = path[:cut_off_last_char_len]
    }
	pathSlice := strings.Split(path, "/")
	return pathSlice
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
