// Package compress provides middleware for request compression.
package compress

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGzipCompression(t *testing.T) {
    handler := http.HandlerFunc(GzipMiddleware(webhook))
    
    srv := httptest.NewServer(handler)
    defer srv.Close()
    
    requestBody := `{
        "request": {
            "type": "SimpleUtterance",
            "command": "sudo do something"
        },
        "version": "1.0"
    }`

    // ожидаемое содержимое тела ответа при успешном запросе
    successBody := `{
        "response": {
            "text": "Извините, я пока ничего не умею"
        },
        "version": "1.0"
    }`
    
    t.Run("sends_gzip", func(t *testing.T) {
        buf := bytes.NewBuffer(nil)
        zb := gzip.NewWriter(buf)
        _, err := zb.Write([]byte(requestBody))
        require.NoError(t, err)
        err = zb.Close()
        require.NoError(t, err)
        
        r := httptest.NewRequest("POST", srv.URL, buf)
        r.RequestURI = ""
        r.Header.Set("Content-Encoding", "gzip")
        r.Header.Set("Accept-Encoding", "")
        
        resp, err := http.DefaultClient.Do(r)
        require.NoError(t, err)
        require.Equal(t, http.StatusOK, resp.StatusCode)
        
        defer resp.Body.Close()
        
        b, err := io.ReadAll(resp.Body)
        require.NoError(t, err)
        require.JSONEq(t, successBody, string(b))
    })

    t.Run("accepts_gzip", func(t *testing.T) {
        buf := bytes.NewBufferString(requestBody)
        r := httptest.NewRequest("POST", srv.URL, buf)
        r.RequestURI = ""
        r.Header.Set("Accept-Encoding", "gzip")
        
        resp, err := http.DefaultClient.Do(r)
        require.NoError(t, err)
        require.Equal(t, http.StatusOK, resp.StatusCode)
        
        defer resp.Body.Close()
        
        zr, err := gzip.NewReader(resp.Body)
        require.NoError(t, err)
        
        b, err := io.ReadAll(zr)
        require.NoError(t, err)
        
        require.JSONEq(t, successBody, string(b))
    })
}

// Helper functions for tests
func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	logger.Log.Debug("decoding request")
	var req Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	text := "Извините, я пока ничего не умею"
	if req.Session.New {
		tz, err := time.LoadLocation(req.Timezone)
		if err != nil {
			logger.Log.Debug("cannot parse timezone")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		now := time.Now().In(tz)
		hour, minute, _ :=  now.Clock()
		text = fmt.Sprintf("Точное время %d часов, %d минут. %s", hour, minute, text)
	}

	resp := Response{
		Response: ResponsePayLoad{
			Text: text,
		},
		Version: "1.0",
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return 
	}
   logger.Log.Debug("sending HTTP 200 response")
}


type Request struct {
	Timezone 	string				`json:"timezone"`
	Request 	SimpleUtterance		`json:"request"`
	Session		Session				`json:"session"`
	Version		string				`json:"version"`
}

type Session	struct {
	New  		bool 				`json:"new"`
}

type SimpleUtterance struct {
	Type 		string				`json:"type"`
	Command		string				`json:"command"`
}

type Response 	struct {
	Response 	ResponsePayLoad		`json:"response"`
	Version		string				`json:"version"`
}

type ResponsePayLoad	struct {
	Text 		string				`json:"text"`
}