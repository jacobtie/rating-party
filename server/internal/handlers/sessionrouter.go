package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jacobtie/rating-party/server/internal/config"
	"github.com/jacobtie/rating-party/server/internal/controllers/session"
	"github.com/jacobtie/rating-party/server/internal/platform/web"
	"github.com/jacobtie/rating-party/server/internal/platform/werrors"
)

type sessionRouter struct {
	controller *session.Controller
}

func registerSessionRoutes(server *web.Service, cfg *config.Config) {
	router := &sessionRouter{
		controller: session.NewController(cfg),
	}
	server.Handle(http.MethodPost, "/signin", router.SignIn)
}

type SignInRequest struct {
	Passcode string `json:"passcode"`
}

func (s *sessionRouter) SignIn(w http.ResponseWriter, r *http.Request) error {
	var request SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return fmt.Errorf("[session.SignIn] failed to decode body with error: %v: %w", err, werrors.ErrBadRequest)
	}
	res, err := s.controller.SignIn(request.Passcode)
	if err != nil {
		return fmt.Errorf("[session.SignIn] failed to sign in: %w", err)
	}
	web.Respond(r.Context(), w, map[string]any{"jwt": res}, http.StatusOK)
	return nil
}
