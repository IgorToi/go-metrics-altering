package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestMekeRequest(t *testing.T) {
	tests := []struct{
		name string 
		path string
	}{
		{
			name: "Simple test #1",
			path: "/update/count/PollCount/5.000000",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			defer srv.Close()
			body, err := MakeRequest(srv.URL + test.path, srv.Client())
			

			assert.Equal(t, 0, len(body))
			assert.NoError(t, err)
		})
	}
}

// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func TestUrlConstructor(t *testing.T) {
	type urlParts struct {
		serverAddress string
		metricType    string
		metricName    string
		valueMetric   float64
	}
	tests := []struct {
        name   string
        urlParts urlParts
        want   string
    }{
		{
			name: "simple test",
			urlParts: urlParts{
				metricType: "count",
				metricName: "PollCount",
				valueMetric: 5,
			},
			want: "http://localhost:8080/update/count/PollCount/5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, flagRunAddr + URLConstructor( protocolScheme, flagRunAddr,  tt.urlParts.metricType, tt.urlParts.metricName, tt.urlParts.valueMetric))
		})
	}
}