package agent

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestSendJSONGauge(t *testing.T) {
	type args struct {
		metricName string
		cfg        *config.ConfigAgent
		value      float64
	}

	successResponse := `"id":"Alloc","type":"gauge","value":1`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var metric models.Metrics

		err := json.NewDecoder(r.Body).Decode(&metric)
		if err != nil {
			log.Println(err)
		}
		assert.Equal(t, "/update/", r.URL.String())
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	cfg := config.ConfigAgent{
		FlagRunAddr: "localhost:8080",
		URL:         server.URL,
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				metricName: "Alloc",
				cfg:        &cfg,
				value:      1.00,
			},
			wantErr: false,
		},
		{
			name: "Matric value not provided",
			args: args{
				metricName: "",
				cfg:        &cfg,
				value:      1.00,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendJSONGauge(tt.args.metricName, tt.args.cfg, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("SendJSONGauge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendURLGauge(t *testing.T) {
	type args struct {
		cfg        *config.ConfigAgent
		value      float64
		metricName string
	}
	successResponse := `"id":"Alloc","type":"gauge","value":1`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//assert.Equal(t, "/update/gauge/Alloc/1.000000", r.URL.String())
		w.Write([]byte(successResponse))
	}))
	defer server.Close()
	cfg := config.ConfigAgent{
		FlagRunAddr: "localhost:8080",
		URL:         server.URL,
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				metricName: "Alloc",
				cfg:        &cfg,
				value:      1.00,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendURLGauge(tt.args.cfg, tt.args.value, tt.args.metricName); (err != nil) != tt.wantErr {
				t.Errorf("SendURLGauge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sendURLCounter(t *testing.T) {
	type args struct {
		cfg        *config.ConfigAgent
		value      int
		metricName string
	}
	successResponse := `"id":"Counter","type":"Counter","value":1`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(successResponse))
	}))
	defer server.Close()
	cfg := config.ConfigAgent{
		FlagRunAddr: "localhost:8080",
		URL:         server.URL,
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				metricName: "Counter",
				cfg:        &cfg,
				value:      1.00,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := sendURLCounter(tt.args.cfg, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("SendURLCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendJSONCounter(t *testing.T) {
	type args struct {
		metricName string
		cfg        *config.ConfigAgent
		value      int
	}

	successResponse := `"id":"Counter","type":"counter","value":1`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var metric models.Metrics

		err := json.NewDecoder(r.Body).Decode(&metric)
		if err != nil {
			log.Println(err)
		}
		assert.Equal(t, "/update/", r.URL.String())
		w.Write([]byte(successResponse))
	}))
	defer server.Close()
	cfg := config.ConfigAgent{
		FlagRunAddr: "localhost:8080",
		URL:         server.URL,
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				metricName: "Counter",
				cfg:        &cfg,
				value:      1.00,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendJSONCounter(tt.args.value, tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("SendJSONCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
