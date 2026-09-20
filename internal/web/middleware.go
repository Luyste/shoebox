package web

import (
	"log"
	"net/http"
	"time"
)

type StatusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *StatusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &StatusRecorder{
			ResponseWriter: w, status: 200,
		}
		start := time.Now()

		next.ServeHTTP(rec, r)

		dur := time.Since(start)
		log.Printf("%v, %v, %v, %v", rec.status, r.Method, r.URL, dur)
	})
}
