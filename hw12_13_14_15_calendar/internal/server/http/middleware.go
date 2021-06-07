package internalhttp

import (
	"net/http"
)

func (s *Server) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		next(w, r)
	}
}
