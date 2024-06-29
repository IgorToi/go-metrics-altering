package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	fmt.Println(req)
	require.NoError(t, err)
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	return resp
}

func TestUpdateHandler(t *testing.T) {
	var m MemStorage
	handler := http.HandlerFunc(m.UpdateHandler)
	srv := httptest.NewServer(handler)
	defer srv.Close()
	sucessBody := `{
		"id": 		"name1",
		"type":	"gauge",
		"value": 	1
	}`
	testCases := []struct{
		name			string
		method			string
		body 			string
		expectedCode	int
		expectedBody	string
	}{
		{
			name: 			"method_get",
			method:			http.MethodGet,
			expectedCode: 	http.StatusMethodNotAllowed,
			expectedBody: 	"",
		},
		{
			name: 			"method_put",
			method:			http.MethodPut,
			expectedCode: 	http.StatusMethodNotAllowed,
			expectedBody: 	"",
		},
		{
			name: 			"method_delete",
			method:			http.MethodDelete,
			expectedCode: 	http.StatusMethodNotAllowed,
			expectedBody: 	"",
		},
		{
			name: 			"method_post_without_body",
			method: 		http.MethodPost,
			expectedCode: 	http.StatusInternalServerError,
			expectedBody: 	"",
		},
		{
			name: 			"method_post_unsupported_type",
			method: 		http.MethodPost,
			body:			`{"id": "name1", "type": "fakeType"}`,
			expectedCode: 	http.StatusUnprocessableEntity,
			expectedBody: 	"",
		},
		{
			name: 			"method_post_success",
			method: 		http.MethodPost,
			body:			`{"id": "name1", "type": "gauge", "value": 1}`,
			expectedCode: 	http.StatusOK,
			expectedBody: 	sucessBody,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL
			if len(tc.body) > 0 {
				req.SetHeader("Content-Type", "application/json")
				req.SetBody(tc.body)
			}
			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(resp.Body()))
			}
		})
	}
}

func TestValueHandler(t *testing.T) {
	m := InitStorage()
	handler := http.HandlerFunc(m.ValueHandler)
	srv := httptest.NewServer(handler)
	defer srv.Close()
	testCases := []struct{
		name			string
		method			string
		body 			string
		expectedCode	int
		expectedBody	string
	}{
		{
			name: 			"method_get",
			method:			http.MethodGet,
			expectedCode: 	http.StatusMethodNotAllowed,
			expectedBody: 	"",
		},
		{
			name: 			"method_put",
			method:			http.MethodPut,
			expectedCode: 	http.StatusMethodNotAllowed,
			expectedBody: 	"",
		},
		{
			name: 			"method_delete",
			method:			http.MethodDelete,
			expectedCode: 	http.StatusMethodNotAllowed,
			expectedBody: 	"",
		},
		{
			name: 			"method_post_without_body",
			method: 		http.MethodPost,
			expectedCode: 	http.StatusInternalServerError,
			expectedBody: 	"",
		},
		{
			name: 			"method_post_unsupported_type",
			method: 		http.MethodPost,
			body:			`{"id": "name1", "type": "fakeType"}`,
			expectedCode: 	http.StatusUnprocessableEntity,
			expectedBody: 	"",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL
			if len(tc.body) > 0 {
				req.SetHeader("Content-Type", "application/json")
				req.SetBody(tc.body)
			}
			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(resp.Body()))
			}
		})
	}
}

func TestInformationHandle(t *testing.T) {
	ts := httptest.NewServer(MetricRouter())
	defer ts.Close()
	type want struct {
		code 		int
		contentType string
	}
	tests := []struct {
		name string
		want want
		url  string
	}{
		{
			name: "Simple test #1",
			want : want{
				code:200,
				contentType: "text/html; charset=utf-8",
			},
			url: "/",
		},
	}
	for _, tt := range tests {
		resp := testRequest(t, ts, "GET", tt.url)
		defer resp.Body.Close()
		assert.Equal(t, tt.want.code, resp.StatusCode)

		assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
	}
}
