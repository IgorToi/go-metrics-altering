package local

import (
	"context"
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
		name       string
		fields     fields
		metricType string
	}{
		{
			name:       "Success",
			fields:     fields{},
			metricType: "counter",
		},
		{
			name:       "Success",
			fields:     fields{},
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
		name      string
		wantError bool
	}{
		{
			name:      "Success",
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

func TestLocalStorage_Update(t *testing.T) {
	m := InitLocalStorage()
	g, c := float64(45), int64(50)

	type args struct {
		ctx         context.Context
		metricType  string
		metricName  string
		metricValue any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success_gauge",
			args: args{
				ctx: context.TODO(),
				metricType: "gauge",
				metricName: "temp_metric1",
				metricValue: &g,
			},
			wantErr: false,
		},
		{
			name: "Success_counter",
			args: args{
				ctx: context.TODO(),
				metricType: "counter",
				metricName: "temp_metric2",
				metricValue: &c,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		_ = m.Update(tt.args.ctx, tt.args.metricType, tt.args.metricName, tt.args.metricValue)

		switch tt.args.metricType {
		case "gauge":
			v, _ :=  tt.args.metricValue.(*float64)
			res := *v

			assert.Equal(t, res, m.Gauge[tt.args.metricName])
		case "counter":
			v, _ :=  tt.args.metricValue.(*int64)
			res := *v

			assert.Equal(t, res, m.Counter[tt.args.metricName])
		}
		
	}
}
