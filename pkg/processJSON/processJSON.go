// Package processjson provides helpful functions for reading and writing JSON requests.
package processjson

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ReadJSON(r *http.Request, dst any) error {
	err := json.NewDecoder(r.Body).Decode(&dst)
	if err != nil {
		return err
	}
	return nil
}

func WriteJSON(rw http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		rw.Header()[key] = value
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(js)
	return nil
}
