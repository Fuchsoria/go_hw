package internalhttp

import (
	"net/http"
	"time"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/app"
)

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}

func loggingMiddleware(next http.HandlerFunc, logger app.Logger) http.HandlerFunc {
	start := time.Now()

	return func(w http.ResponseWriter, r *http.Request) {
		o := &responseObserver{ResponseWriter: w}

		next(o, r)

		logger.Info("",
			"ip", r.RemoteAddr,
			"date", time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			"method", r.Method,
			"path", r.URL.Path,
			"http", r.Proto,
			"code", o.status,
			"latency", time.Since(start),
			"useragent", r.UserAgent())
	}
}
