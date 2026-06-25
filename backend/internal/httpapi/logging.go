package httpapi

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *loggingResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *loggingResponseWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(p)
	w.bytes += n
	return n, err
}

func loggingMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lrw, r)
		status := lrw.status
		if status == 0 {
			status = http.StatusOK
		}
		logger.Printf(
			"event=http_request method=%s path=%s status=%d bytes=%d duration=%s",
			r.Method,
			requestPath(r),
			status,
			lrw.bytes,
			time.Since(start).Truncate(time.Millisecond),
		)
	})
}

func requestPath(r *http.Request) string {
	if strings.TrimSpace(r.URL.RawQuery) == "" {
		return r.URL.Path
	}
	return r.URL.Path + "?" + r.URL.RawQuery
}
