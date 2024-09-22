package timeout

import (
	"context"
	"net/http"
	"time"

	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
)

func Timeout(timeout time.Duration, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		r = r.WithContext(ctx)

		processDone := make(chan bool)
		go func() {
			next.ServeHTTP(w, r)
			processDone <- true
		}()

		select {
		case <-ctx.Done():
			logger.Log.Info("HTTP Request timed out")
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("timed out"))
		case <-processDone:
		}
	})
}
