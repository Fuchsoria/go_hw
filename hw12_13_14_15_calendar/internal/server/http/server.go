package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	app    app.App
	server *http.Server
}

type Handler struct {
	app app.App
}

type Message struct {
	Message string `json:"message"`
}

type Results struct {
	Results []storage.Event `json:"results"`
}

func NewServer(app *app.App, address string, port string) *Server {
	h := Handler{*app}
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", loggingMiddleware(h.Hello, app.GetLogger()))
	mux.HandleFunc("/events/create", loggingMiddleware(h.Create, app.GetLogger()))
	mux.HandleFunc("/events/update", loggingMiddleware(h.Update, app.GetLogger()))
	mux.HandleFunc("/events/delete", loggingMiddleware(h.Delete, app.GetLogger()))
	mux.HandleFunc("/events/daily", loggingMiddleware(h.Daily, app.GetLogger()))
	mux.HandleFunc("/events/weekly", loggingMiddleware(h.Weekly, app.GetLogger()))
	mux.HandleFunc("/events/month", loggingMiddleware(h.Month, app.GetLogger()))

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
		return fmt.Errorf("cannot start http server, %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("cannot shutdown http server, %w", err)
	}

	return nil
}

func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("ok")
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	bData := struct {
		Title string `json:"title"`
		Date  int    `json:"date"`
	}{}

	err := h.getBodyData(w, r, &bData)
	if err != nil {
		h.send(w, http.StatusInternalServerError, Message{err.Error()})
	}

	err = h.app.CreateEvent(bData.Title, int64(bData.Date))
	if err != nil {
		h.send(w, http.StatusInternalServerError, Message{"Server error"})

		return
	}

	h.send(w, http.StatusCreated, Message{"Created new"})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	event := storage.Event{}

	err := h.getBodyData(w, r, &event)
	if err != nil {
		h.send(w, http.StatusInternalServerError, Message{err.Error()})
	}

	err = h.app.UpdateEvent(event)
	if err != nil {
		h.send(w, http.StatusInternalServerError, Message{err.Error()})

		return
	}

	h.send(w, http.StatusOK, Message{"Updated"})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	bData := struct {
		ID string `json:"id"`
	}{}

	err := h.getBodyData(w, r, &bData)
	if err != nil {
		h.send(w, http.StatusInternalServerError, Message{err.Error()})
	}

	err = h.app.RemoveEvent(bData.ID)
	if err != nil {
		h.send(w, http.StatusInternalServerError, Message{err.Error()})

		return
	}

	h.send(w, http.StatusOK, Message{"Deleted"})
}

func (h *Handler) Daily(w http.ResponseWriter, r *http.Request) {
	h.getEventsByDate(w, r, "day")
}

func (h *Handler) Weekly(w http.ResponseWriter, r *http.Request) {
	h.getEventsByDate(w, r, "week")
}

func (h *Handler) Month(w http.ResponseWriter, r *http.Request) {
	h.getEventsByDate(w, r, "month")
}

func (h *Handler) getEventsByDate(w http.ResponseWriter, r *http.Request, dateType string) {
	bData := struct {
		Date int `json:"date"`
	}{}

	err := h.getBodyData(w, r, &bData)
	if err != nil {
		h.send(w, http.StatusInternalServerError, Message{err.Error()})
	}

	var (
		events []storage.Event
		time   time.Time = time.Unix(int64(bData.Date), 0)
	)

	switch dateType {
	case "day":
		events, err = h.app.DailyEvents(time)
	case "week":
		events, err = h.app.WeeklyEvents(time)
	case "month":
		events, err = h.app.MonthEvents(time)
	}
	if err != nil {
		h.send(w, http.StatusInternalServerError, Message{err.Error()})
	}

	h.send(w, http.StatusOK, Results{events})
}

func (h *Handler) getBodyData(w http.ResponseWriter, r *http.Request, bData interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.send(w, http.StatusBadRequest, Message{"Bad request"})

		return fmt.Errorf("cannot read body, %w", err)
	}

	if err = json.Unmarshal(b, &bData); err != nil {
		h.send(w, http.StatusBadRequest, Message{"Bad request"})

		return fmt.Errorf("cannot unmarshall body, %w", err)
	}

	return nil
}

func (h *Handler) send(w http.ResponseWriter, status int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(resp)
}
