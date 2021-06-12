package internalhttp

import (
	"net/http"
	"time"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/app"
)

func loggingMiddleware(next http.HandlerFunc, logger app.Logger) http.HandlerFunc {
	start := time.Now()

	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)

		logger.Info("", "ip", r.RemoteAddr, "date", time.Now().Format("02/Jan/2006:15:04:05 -0700"), "method", r.Method, "path", r.URL.Path, "http", r.Proto, "code", 200, "latency", time.Since(start), "useragent", r.UserAgent())
	}
}
