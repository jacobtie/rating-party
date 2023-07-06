package web

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jacobtie/rating-party/server/internal/platform/contextvalue"
	"github.com/julienschmidt/httprouter"
)

type Service struct {
	router           *httprouter.Router
	globalMiddleware []Middleware
}

func NewService(mw ...Middleware) *Service {
	s := &Service{
		router:           httprouter.New(),
		globalMiddleware: mw,
	}
	s.router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			header.Set("Access-Control-Allow-Headers", "*")
		}
		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})
	return s
}

type Handler func(http.ResponseWriter, *http.Request) error

func (s *Service) Handle(verb, path string, handler Handler, mw ...Middleware) {
	wrappedHandler := wrapMiddleware(wrapMiddleware(handler, mw), s.globalMiddleware)
	h := func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		header.Set("Access-Control-Allow-Headers", "*")
		ctx := r.Context()
		requestID := uuid.New().String()
		v := &contextvalue.Values{
			RequestID:    requestID,
			RequestStart: time.Now(),
		}
		ctx = context.WithValue(ctx, contextvalue.KeyValues, v)
		if err := wrappedHandler(w, r.WithContext(ctx)); err != nil {
			HandleError(r.Context(), w, err)
			return
		}
	}
	s.router.HandlerFunc(verb, path, h)
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
