package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	_ "net/http/pprof"
	"testing"

	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/igortoigildin/go-metrics-altering/internal/server/api/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_updates(t *testing.T) {
	gaugeValue := float64(1.5)
	counterValue := int64(5)

	type metric struct {
		ID    string `json:"id"`
		MType string `json:"type"`
		Value any    `json:"value,omitempty"`
		Delta any    `json:"delta,omitempty"`
	}

	input := []metric{
		{
			ID:    "first_metric",
			MType: "gauge",
			Value: &gaugeValue,
		},
		{
			ID:    "second_metric",
			MType: "incorrect_type",
			Value: &gaugeValue,
		},
		{
			ID:    "third_metric",
			MType: "gauge",
			Value: "adfad",
		},
		{
			ID:    "4th_metric",
			MType: "gauge",
			Value: &gaugeValue,
		},
		{
			ID:    "5th_metric",
			MType: "counter",
			Delta: &counterValue,
		},
	}

	tests := []struct {
		name           string
		arguments      []metric // udpated metrics which need to be saved
		respError      string
		mockError      error
		respStatusCode int
		inputIndex     []int
	}{
		{
			name:           "Success",
			inputIndex:     []int{0, 4},
			arguments:      input,
			respStatusCode: 200,
		},
		{
			name:           "Unsupported metric type",
			inputIndex:     []int{1},
			arguments:      input,
			respError:      "usupported request type",
			respStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "Usupported metric value",
			inputIndex:     []int{2},
			arguments:      input,
			respError:      "Usupported metric value",
			respStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Storage error",
			inputIndex:     []int{3},
			arguments:      input,
			respStatusCode: http.StatusInternalServerError,
			mockError:      errors.New("unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewStorage(t)

			if tt.respError == "" && tt.mockError == nil {
				repo.On("Update", mock.Anything, tt.arguments[0].MType, tt.arguments[0].ID, tt.arguments[0].Value).Return(nil).Times(1)
				repo.On("Update", mock.Anything, tt.arguments[4].MType, tt.arguments[4].ID, tt.arguments[4].Delta).Return(nil).Times(1)
			}

			if tt.mockError != nil {
				repo.On("Update", mock.Anything, tt.arguments[3].MType, tt.arguments[3].ID, tt.arguments[3].Value).Return(tt.mockError).Times(1)
			}

			var metrics []metric
			for _, v := range tt.inputIndex {
				metrics = append(metrics, tt.arguments[v])
			}
			js, _ := json.Marshal(metrics)

			handler := updates(repo)
			req, err := http.NewRequest(http.MethodPost, "/updates/", bytes.NewReader([]byte(js)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.respStatusCode, rr.Code)
		})
	}
}

func Test_updateMetric(t *testing.T) {
	gaugeValue := float64(1.5)
	counterValue := int64(5)

	type metric struct {
		ID    string `json:"id"`
		MType string `json:"type"`
		Value any    `json:"value,omitempty"`
		Delta any    `json:"delta,omitempty"`
	}

	input := []metric{
		{
			ID:    "first_metric",
			MType: "gauge",
			Value: &gaugeValue,
		},
		{
			ID:    "second_metric",
			MType: "incorrect_type",
			Value: &gaugeValue,
		},
		{
			ID:    "third_metric",
			MType: "gauge",
			Value: "adfad",
		},
		{
			ID:    "4th_metric",
			MType: "gauge",
			Value: &gaugeValue,
		},
		{
			ID:    "5th_metric",
			MType: "counter",
			Delta: &counterValue,
		},
	}

	tests := []struct {
		name           string
		arguments      []metric // udpated metric which need to be saved
		respError      string
		mockError      error
		respStatusCode int
		inputIndex     int
		method         string
	}{
		{
			name:           "Success",
			inputIndex:     0,
			arguments:      input,
			respStatusCode: 200,
			method:         http.MethodPost,
		},
		{
			name:           "Unsupported metric type",
			inputIndex:     1,
			arguments:      input,
			respError:      "usupported request type",
			respStatusCode: http.StatusUnprocessableEntity,
			method:         http.MethodPost,
		},
		{
			name:           "Usupported metric value",
			inputIndex:     2,
			arguments:      input,
			respError:      "Usupported metric value",
			respStatusCode: http.StatusBadRequest,
			method:         http.MethodPost,
		},
		{
			name:           "Storage error",
			inputIndex:     3,
			arguments:      input,
			respStatusCode: http.StatusInternalServerError,
			mockError:      errors.New("unexpected error"),
			method:         http.MethodPost,
		},
		{
			name:           "Incorrect method",
			inputIndex:     3,
			arguments:      input,
			respError:      "Method not allowed",
			respStatusCode: http.StatusMethodNotAllowed,
			method:         http.MethodGet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewStorage(t)

			if tt.respError == "" && tt.mockError == nil {
				repo.On("Update", mock.Anything, tt.arguments[0].MType, tt.arguments[0].ID, tt.arguments[0].Value).Return(nil).Times(1)
			}

			if tt.mockError != nil {
				repo.On("Update", mock.Anything, tt.arguments[3].MType, tt.arguments[3].ID, tt.arguments[3].Value).Return(tt.mockError).Times(1)
			}

			var metrics []metric
			metrics = append(metrics, tt.arguments[tt.inputIndex])
			js, _ := json.Marshal(metrics)

			handler := updates(repo)
			req, err := http.NewRequest(tt.method, "/update/", bytes.NewReader([]byte(js)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.respStatusCode, rr.Code)
		})
	}
}

func Test_getMetric(t *testing.T) {
	gaugeValue := float64(10)
	type metric struct {
		ID    string `json:"id"`
		MType string `json:"type"`
	}

	input := []metric{
		{
			ID:    "first_metric",
			MType: "gauge",
		},
		{
			ID:    "second_metric",
			MType: "incorrect_type",
		},
		{
			ID:    "third_metric",
			MType: "counter",
		},
	}

	tests := []struct {
		name           string
		arguments      []metric // metrics which need to be requested
		respError      string
		mockError      error
		respStatusCode int
		inputIndex     int
		response       models.Metrics
	}{
		{
			name:           "Success with gauge",
			inputIndex:     0,
			arguments:      input,
			respStatusCode: 200,
			response: models.Metrics{
				ID:    "first_metric",
				MType: "gauge",
				Value: &gaugeValue,
			},
		},
		{
			name:           "Unsupported metric type",
			inputIndex:     1,
			arguments:      input,
			respError:      "usupported request type",
			respStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "Success with counter",
			inputIndex:     2,
			arguments:      input,
			respStatusCode: http.StatusOK,
			response: models.Metrics{
				ID:    "next_metric",
				MType: "gauge",
				Value: &gaugeValue,
			},
		},
		{
			name:           "Storage error",
			inputIndex:     0,
			arguments:      input,
			respStatusCode: http.StatusInternalServerError,
			mockError:      errors.New("unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewStorage(t)

			if tt.respError == "" && tt.mockError == nil {
				repo.On("Get", mock.Anything, tt.arguments[0].MType, tt.arguments[0].ID).Return(tt.response, nil).Maybe()
				repo.On("Get", mock.Anything, tt.arguments[2].MType, tt.arguments[2].ID).Return(tt.response, nil).Maybe()
			}

			if tt.mockError != nil {
				repo.On("Get", mock.Anything, tt.arguments[0].MType, tt.arguments[0].ID).Return(models.Metrics{}, tt.mockError).Times(1)
			}

			js, _ := json.Marshal(tt.arguments[tt.inputIndex])

			handler := getMetric(repo)
			req, err := http.NewRequest(http.MethodPost, "/value/", bytes.NewReader([]byte(js)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.respStatusCode, rr.Code)
		})
	}
}

func Test_ping(t *testing.T) {
	tests := []struct {
		name           string
		response       map[string]interface{}
		respError      string
		mockError      error
		respStatusCode int
		method         string
	}{
		{
			name:           "Success",
			respStatusCode: http.StatusOK,
			method:         http.MethodGet,
		},
		{
			name:           "Wrong method",
			respStatusCode: http.StatusMethodNotAllowed,
			method:         http.MethodPost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewStorage(t)

			if tt.respError == "" && tt.mockError == nil {
				repo.On("Ping", context.Background()).Return(nil).Maybe()
			}

			handler := ping(repo)

			req, err := http.NewRequest(tt.method, "/ping", bytes.NewReader([]byte(nil)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.respStatusCode, rr.Code)
		})
	}
}

// "/update/{metricType}/{metricName}/{metricValue}"
// func Test_updatePathHandler(t *testing.T) {
// 	repo := mocks.NewStorage(t)
// 	handler := http.HandlerFunc(updatePathHandler(repo))
//     srv := httptest.NewServer(handler)
//     defer  srv.Close()

// 	type metric struct {
// 		ID    string `json:"id"`
// 		MType string `json:"type"`
// 		Value float64    `json:"value,omitempty"`
// 		Delta any    `json:"delta,omitempty"`
// 	}

// 	input := []metric{
// 		{
// 			ID:    "first_metric",
// 			MType: "gauge",
// 			Value: float64(5),
// 		},
// 		// {
// 		// 	ID:    "second_metric",
// 		// 	MType: "incorrect_type",
// 		// 	Value: &gaugeValue,
// 		// },
// 		// {
// 		// 	ID:    "third_metric",
// 		// 	MType: "gauge",
// 		// 	Value: "adfad",
// 		// },
// 		// {
// 		// 	ID:    "4th_metric",
// 		// 	MType: "gauge",
// 		// 	Value: &gaugeValue,
// 		// },
// 		// {
// 		// 	ID:    "5th_metric",
// 		// 	MType: "counter",
// 		// 	Delta: &counterValue,
// 		// },
// 	}

// 	tests := []struct {
// 		name           string
// 		arguments      []metric // udpated metric which need to be saved
// 		respError      string
// 		mockError      error
// 		respStatusCode int
// 		inputIndex     int
// 		method         string
// 	}{
// 		{
// 			name:           "Success",
// 			inputIndex:     0,
// 			arguments:      input,
// 			respStatusCode: 200,
// 			method:         http.MethodPost,
// 		},
// 		// {
// 		// 	name:           "Unsupported metric type",
// 		// 	inputIndex:     1,
// 		// 	arguments:      input,
// 		// 	respError:      "usupported request type",
// 		// 	respStatusCode: http.StatusUnprocessableEntity,
// 		// 	method:         http.MethodPost,
// 		// },
// 		// {
// 		// 	name:           "Usupported metric value",
// 		// 	inputIndex:     2,
// 		// 	arguments:      input,
// 		// 	respError:      "Usupported metric value",
// 		// 	respStatusCode: http.StatusBadRequest,
// 		// 	method:         http.MethodPost,
// 		// },
// 		// {
// 		// 	name:           "Storage error",
// 		// 	inputIndex:     3,
// 		// 	arguments:      input,
// 		// 	respStatusCode: http.StatusInternalServerError,
// 		// 	mockError:      errors.New("unexpected error"),
// 		// 	method:         http.MethodPost,
// 		// },
// 		// {
// 		// 	name:           "Incorrect method",
// 		// 	inputIndex:     3,
// 		// 	arguments:      input,
// 		// 	respError:      "Method not allowed",
// 		// 	respStatusCode: http.StatusMethodNotAllowed,
// 		// 	method:         http.MethodGet,
// 		// },
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			repo := mocks.NewStorage(t)

// 			if tt.respError == "" && tt.mockError == nil {
// 				repo.On("Update", mock.Anything, tt.arguments[0].MType, tt.arguments[0].ID, tt.arguments[0].Value).Return(nil).Times(1)
// 			}

// 			// if tt.mockError != nil {
// 			// 	repo.On("Update", mock.Anything, tt.arguments[3].MType, tt.arguments[3].ID, tt.arguments[3].Value).Return(tt.mockError).Times(1)
// 			// }

// 			//agent := resty.New()
// 			req := resty.New().R().SetPathParams(map[string]string{
// 				"metricType":  "gauge",
// 				"metricName":  "name",
// 				"metricValue": strconv.FormatFloat(tt.arguments[0].Value, 'f', 6, 64),
// 			}).SetHeader("Content-Type", "text/plain")


// 			resp, err := req.Post(srv.URL + "/update/{metricType}/{metricName}/{metricValue}")
// 			fmt.Println(resp.StatusCode())
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 			fmt.Println(req.URL)

// 			assert.Equal(t, tt.respStatusCode, resp.StatusCode(), "Response code didn't match expected")
// 			require.NoError(t, err)
// 		})
// 	}
// }
