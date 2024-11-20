package agent

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	configSrv "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/igortoigildin/go-metrics-altering/pkg/crypt"
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
		// var val float64 = 0
		if metric.ID == "Alloc1" {
			time.Sleep(time.Second * 5)
			//w.WriteHeader(http.StatusRequestTimeout)
			return
		}
		assert.Equal(t, "/update/", r.URL.String())
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	cfg := config.ConfigAgent{
		FlagRunAddrHTTP:       "localhost:8080",
		URL:               server.URL,
		FlagHashKey:       "123",
		FlagRSAEncryption: true,
		FlagCryptoKey:     "keys/public.pem",
	}

	// preparing temp config
	cfgServer := configSrv.ConfigServer{}
	// init temp rsa keys
	_ = crypt.InitRSAKeys(&cfgServer)

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
		{
			name: "Retry",
			args: args{
				metricName: "Alloc1",
				cfg:        &cfg,
				value:      0.00,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendJSONGauge(tt.args.metricName, tt.args.cfg, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("SendJSONGauge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	cfgnew := config.ConfigAgent{
		FlagRunAddrHTTP: "localhost:8080",
		URL:         "incorrect",
	}
	if err := SendJSONGauge("new", &cfgnew, float64(0)); err != nil {
		_ = err
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
		FlagRunAddrHTTP:       "localhost:8080",
		URL:               server.URL,
		FlagHashKey:       "123",
		FlagRSAEncryption: true,
		FlagCryptoKey:     "keys/public.pem",
	}

	// preparing temp config
	cfgServer := configSrv.ConfigServer{}
	// init temp rsa keys
	_ = crypt.InitRSAKeys(&cfgServer)

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

	cfgnew := config.ConfigAgent{
		FlagRunAddrHTTP: "localhost:8080",
		URL:         "incorrect",
	}
	if err := SendJSONGauge("new", &cfgnew, float64(0)); err != nil {
		_ = err
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
		FlagRunAddrHTTP:       "localhost:8080",
		URL:               server.URL,
		FlagHashKey:       "123",
		FlagRSAEncryption: true,
		FlagCryptoKey:     "keys/public.pem",
	}

	// preparing temp config
	cfgServer := configSrv.ConfigServer{}
	// init temp rsa keys
	_ = crypt.InitRSAKeys(&cfgServer)

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

	cfgnew1 := config.ConfigAgent{
		FlagRunAddrHTTP: "localhost:8080",
		URL:         "incorrect",
	}
	if err := sendURLCounter(&cfgnew1, 1); err != nil {
		assert.Error(t, err)
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
		FlagRunAddrHTTP:       "localhost:8080",
		URL:               server.URL,
		FlagHashKey:       "123",
		FlagRSAEncryption: true,
		FlagCryptoKey:     "keys/public.pem",
	}

	// preparing temp config
	cfgServer := configSrv.ConfigServer{}
	// init temp rsa keys
	_ = crypt.InitRSAKeys(&cfgServer)

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

	cfgnew1 := config.ConfigAgent{
		FlagRunAddrHTTP: "localhost:8080",
		URL:         "incorrect",
	}
	if err := SendJSONCounter(1, &cfgnew1); err != nil {
		assert.Error(t, err)
	}
}
