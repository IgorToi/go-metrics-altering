package local

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalStorage_SetStrategy(t *testing.T) {
	type fields struct {
		rm       sync.RWMutex
		Gauge    map[string]float64
		Counter  map[string]int64
		strategy Strategy
	}

	tests := []struct {
		name   string
		fields fields
		metricType   string
	}{
		{
			name: "Success",
			fields: fields{},
			metricType: "counter",
		},
		{
			name: "Success",
			fields: fields{},
			metricType: "gauge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &LocalStorage{
				rm:       tt.fields.rm,
				Gauge:    tt.fields.Gauge,
				Counter:  tt.fields.Counter,
				strategy: tt.fields.strategy,
			}
			m.SetStrategy(tt.metricType)

			assert.True(t, m.strategy != nil)
		})
	}
}

func TestInitLocalStorage(t *testing.T) {
	tests := []struct {
		name string
		wantError bool
	}{
		{
			name: "Success",
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := InitLocalStorage()
			assert.True(t, a.Counter != nil)
			assert.True(t, a.Gauge != nil)
		})
	}
}
