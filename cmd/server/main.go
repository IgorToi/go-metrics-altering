package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/go-chi/chi"
)

const htmlTemplate = `
{{range $index, $element := .}}
{{ $index }}: {{ $element }}
{{end}}
`
type MemStorage struct {
	gauge     map[string]float64
	counter   map[string]int64
}

var Memory = &MemStorage{
	gauge: make(map[string]float64),
	counter: make(map[string]int64),
}

func (m *MemStorage) UpdateGaugeMetric(metricName string, metricValue float64) {
	m.gauge[metricName] = metricValue
}

func (m *MemStorage) UpdateCounterMetric(metricName string, metricValue int64) {
	m.counter[metricName] += metricValue
}

func (m *MemStorage) GetGaugeMetricFromMemory(metricName string) float64 {
	return m.gauge[metricName]
}

func (m *MemStorage) GetCountMetricFromMemory(metricName string) int64 {
	return m.counter[metricName]
}

func (m *MemStorage) CheckIfGaugeMetricPresent(metricName string) bool {
	_, ok := m.gauge[metricName]
	return ok
}

func (m *MemStorage) CheckIfCountMetricPresent(metricName string) bool {
	_, ok := m.counter[metricName]
	return ok
}


func UpdateHandle(rw http.ResponseWriter, r *http.Request) { 
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	if metricName ==  "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	metricValue := chi.URLParam(r, "metricValue")

	switch metricType {
	case "gauge":
		metricValueConverted, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		Memory.UpdateGaugeMetric(metricName, metricValueConverted)
	case "counter":
		metricValueConverted, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		Memory.UpdateCounterMetric(metricName, metricValueConverted)

	default:
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func ValueHandle(rw http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	switch metricType {
	case "gauge":
		if !Memory.CheckIfGaugeMetricPresent(metricName) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.Write([]byte(strconv.FormatFloat(Memory.GetGaugeMetricFromMemory(metricName),'f', -1, 64)))
	case "counter":
		if !Memory.CheckIfCountMetricPresent(metricName) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.Write([]byte(strconv.FormatInt(Memory.GetCountMetricFromMemory(metricName), 10)))
	default:
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func InformationHandle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.New("t")
    t, err := t.Parse(htmlTemplate)
    if err != nil {
        panic(err)
    }
	
    err = t.Execute(rw, convertToSingleMap(Memory.gauge, Memory.counter))
    if err != nil {
        panic(err)
    }
}

func MetricRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", UpdateHandle)
	r.Get("/value/{metricType}/{metricName}", ValueHandle)
	r.Get("/", InformationHandle)
	return r
}


func main() {	
	parseFlags()

	fmt.Println("Running server on", flagRunAddr)
	http.ListenAndServe(flagRunAddr, MetricRouter())
}

func convertToSingleMap(a map[string]float64, b map[string]int64) map[string]interface{} {
	c := make(map[string]interface{})
	for i, v := range a {
		c[i] = v
	}
	for j, l := range b {
		c[j] = l
	}
	return c
}