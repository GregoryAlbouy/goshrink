package httplog

import (
	"log"
	"net/http"
)

// ResponseWriter is a wrapper around http.ResponseWriter that provides
// extra information on the request.
type ResponseWriter struct {
	http.ResponseWriter
	Status int // Records the status code of the request.
}

func (lr *ResponseWriter) WriteHeader(statusCode int) {
	lr.Status = statusCode
	lr.ResponseWriter.WriteHeader(statusCode)
}

// RequestLogger adds logging to the given http.Handler.
func RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lr := &ResponseWriter{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}

		c := statusColor(lr.Status)

		defer log.Printf("%s %s -> %s %s", r.Method, r.URL.String(), c(lr.Status), c(http.StatusText(lr.Status)))
		h.ServeHTTP(lr, r)
	})
}
