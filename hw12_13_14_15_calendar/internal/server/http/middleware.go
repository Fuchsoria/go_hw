package internalhttp

import (
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	start := time.Now()

	return func(w http.ResponseWriter, r *http.Request) {
		logger := s.app.GetLogger()

		next(w, r)

		logger.Info("", "ip", r.RemoteAddr, "date", time.Now().Format("02/Jan/2006:15:04:05 -0700"), "method", r.Method, "path", r.URL.Path, "http", r.Proto, "code", 200, "latency", time.Since(start), "useragent", r.UserAgent())
	}
}
