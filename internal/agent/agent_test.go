package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestSendMetric(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {}))
	defer server.Close()

	tests := []struct {
		name         string
		method       string
		responseCode int
		metricType   string
		metricName   string
		metricValue  string
	}{
		{
			name:         "Simple test #1",
			method:       "POST",
			responseCode: http.StatusOK,
			metricType:   "counter",
			metricName:   "PollCount",
			metricValue:  "527",
		},
		{
			name:         "Simple test #2",
			method:       "POST",
			responseCode: http.StatusOK,
			metricType:   "gauge",
			metricName:   "Alloc",
			metricValue:  "500",
		},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			agent := resty.New()
			req := agent.R()
			req.SetPathParams(map[string]string{
				"metricType":  tt.metricType,
				"metricName":  tt.metricName,
				"metricValue": tt.metricValue,
			}).SetHeader("Content-Type", "text/plain")
			req.URL = server.URL
			resp, err := SendMetric(req.URL, tt.metricType, tt.metricName, tt.metricValue, req)

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tt.responseCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}
