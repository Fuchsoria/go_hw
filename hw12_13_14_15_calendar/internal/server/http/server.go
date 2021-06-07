package internalhttp

import (
	"context"
	"log"
	"net"
	"net/http"
)

type Server struct { // TODO
	app Application
}

type Application interface {
}

func NewServer(app Application) *Server {
	return &Server{app}
}

func (s *Server) Start(ctx context.Context, address string, port string) error {
	handler := s
	mux := http.NewServeMux()
	mux.HandleFunc("/search", s.loggingMiddleware(handler.Search))
	mux.HandleFunc("/add", s.loggingMiddleware(handler.AddItem))

	server := &http.Server{
		Addr:    net.JoinHostPort(address, port),
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())

	http.ListenAndServe(net.JoinHostPort(address, port), nil)
	// TODO
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Server) Search(w http.ResponseWriter, r *http.Request) {
	// ...
}
func (s *Server) AddItem(w http.ResponseWriter, r *http.Request) {
	// ...
}
