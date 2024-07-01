package server

import (
	"encoding/json"
	"os"
)

type Metrics struct {
	Alloc 			float64
	BuckHashSys		float64
	Frees			float64
	GCCPUFraction	float64
	GCSys			float64
	HeapAlloc		float64
	HeapIdle		float64
	HeapInuse		float64
	HeapObjects		float64
	HeapReleased	float64
	HeapSys			float64
	LastGC			float64
	Lookups			float64
	MCacheInuse		float64
	MSpanSys 		float64
	MCacheSys		float64
	MSpanInuse		float64
	Mallocs			float64
	NextGC			float64
	NumForcedGC		float64
	NumGC			float64
	OtherSys		float64
	PauseTotalNs	float64
	StackInuse		float64
	StackSys		float64
	Sys				float64
	TotalAlloc		float64
	RandomValue		float64

	PollCount		int
}





















func (metrics Metrics) Save(fname string)  error {
	data, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fname, data, 0606)
}

func (metrics *Metrics) Load(fname string) error {
	data, err := os.ReadFile(fname)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, metrics); err != nil {
		return err
	}
	return nil
}

