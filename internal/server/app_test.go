package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/internal/mocks"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := mocks.NewMockStorage(ctrl)
	testValue := float64(10.0)
	metric := models.Metrics{
		ID:    "Alloc",
		MType: config.GaugeType,
		Delta: nil,
		Value: &testValue,
	}
	s.EXPECT().Get(gomock.Any(), config.GaugeType, "Alloc").Return(metric, nil)
	s.EXPECT().Exist(gomock.Any(), config.GaugeType, "Alloc").Return(true)

	appInstance := newApp(s)
	handler := http.HandlerFunc(appInstance.getMetric)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode string
		expectedBody string
	}{
		{
			name:         "method_post_success",
			method:       http.MethodPost,
			body:         `{"id":"Alloc", "type":"gauge"}`,
			expectedCode: "200 OK",
			expectedBody: `{"id": "Alloc", "type": "gauge", "delta": "", "value": "10.0"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			var jsonStr = []byte(tc.body)

			NewReq, err := http.NewRequest("POST", srv.URL+"/value/", bytes.NewReader(jsonStr))
			NewReq.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			resp, _ := client.Do(NewReq)

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.expectedCode, resp.Status, "Response code didn't match expected")
			if tc.expectedBody != "" {

				_, _ = io.ReadAll(resp.Body)
				defer resp.Body.Close()

				// to be updated after code refactoring

				//assert.JSONEq(t, tc.expectedBody, string(body))
			}
		})
	}
}
