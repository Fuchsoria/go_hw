package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	app    app.App
	server *http.Server
}

type Handler struct{}

func NewServer(app *app.App, address string, port string) *Server {
	h := Handler{}
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", loggingMiddleware(h.Hello, app.GetLogger()))

	server := &http.Server{
		Addr:         net.JoinHostPort(address, port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{*app, server}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("can not start http server, %w", err)
	}

	<-ctx.Done()

	if err := s.Stop(ctx); err != nil {
		return fmt.Errorf("can not stop http server, %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("can not shutdown http server, %w", err)
	}

	return nil
}

func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("ok")
}
