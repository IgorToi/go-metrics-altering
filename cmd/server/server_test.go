package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp
}


func TestRouter(t *testing.T) {
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
			name: "positive test #1",
			want: want{
				code: 200,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/counter/someMetric/527",
		},
		{
			name: "negative test #2",
			want: want{
				code: 404,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/counter/527",
		},
		{
			name: "negative test #3",
			want: want{
				code: 400,
				contentType: "",
			},
			url: "/update/counter/someMetric/527_e",
		},
		{
			name: "negative test #4",
			want: want{
				code: 400,
				contentType: "",
			},
			url: "/update/counter1/someMetric/527",
		},
	}
	for _, tt := range tests {
		resp := testRequest(t, ts, "POST", tt.url)
		assert.Equal(t, tt.want.code, resp.StatusCode)

		assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
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
		assert.Equal(t, tt.want.code, resp.StatusCode)

		assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
	}
}

func TestValueHandle(t *testing.T) {
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
			name: "positive test #1",
			want: want{
				code: 200,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/value/gauge/Alloc",
		},
		{
			name: "negative test #2",
			want: want{
				code: 404,
				contentType: "",
			},
			url: "/value/gauge/A",
		},
	}
	for _, tt := range tests {
		Memory.gauge["Alloc"] = 45
		resp := testRequest(t, ts, "GET", tt.url)
		assert.Equal(t, tt.want.code, resp.StatusCode)

		assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
	}
}