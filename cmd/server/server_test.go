package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReqeustHandler(t *testing.T) {
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
				contentType: "text/plain",
			},
			url: "http://localhost:8080/update/counter/someMetric/527",
		},
		{
			name: "negative test #2",
			want: want{
				code: 404,
				contentType: "",
			},
			url: "http://localhost:8080/update/counter/527",
		},
		{
			name: "negative test #3",
			want: want{
				code: 400,
				contentType: "",
			},
			url: "http://localhost:8080/update/string/someMetric/527",
		},
		{
			name: "negative test #4",
			want: want{
				code: 400,
				contentType: "",
			},
			url: "http://localhost:8080/update/counter/someMetric/52e7",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()
			ReqeustHandler(w, request)

			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)

			defer res.Body.Close()
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestPathCleaner(t *testing.T) {
	type want struct {
		path []string
	}
	tests := []struct {
		name string
		url string
		want want
	}{
		{
			name: "simple test #1",
			want: want{
				path: []string{"update","counter", "someMetric", "527"},
			},
			url: "/update/counter/someMetric/527",
		},
		{
			name: "simple test #2",
			want: want{
				path: []string{"update","counter", "someMetric", "527"},
			},
			url: "update/counter/someMetric/527",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want.path, pathCleaner(tt.url))
		})
	}
}
