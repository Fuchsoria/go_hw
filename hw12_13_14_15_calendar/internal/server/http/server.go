package internalhttp

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/app"
)

type Server struct { // TODO
	app    Application
	server *http.Server
}

type Application interface {
	GetLogger() app.Logger
}

func NewServer(app Application) *Server {
	return &Server{app, nil}
}

func (s *Server) Start(ctx context.Context, address string, port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", s.loggingMiddleware(s.Hello))

	s.server = &http.Server{
		Addr:         net.JoinHostPort(address, port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(s.server.ListenAndServe())

	http.ListenAndServe(net.JoinHostPort(address, port), nil)
	// TODO
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Server) Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("ok")
}
