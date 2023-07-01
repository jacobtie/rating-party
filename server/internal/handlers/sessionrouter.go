package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jacobtie/rating-party/server/internal/config"
	"github.com/jacobtie/rating-party/server/internal/controllers/session"
	"github.com/jacobtie/rating-party/server/internal/platform/db"
	"github.com/jacobtie/rating-party/server/internal/platform/web"
	"github.com/jacobtie/rating-party/server/internal/platform/werrors"
)

type sessionRouter struct {
	controller *session.Controller
}

func registerSessionRoutes(server *web.Service, cfg *config.Config, db *db.DB) {
	router := &sessionRouter{
		controller: session.NewController(cfg, db),
	}
	server.Handle(http.MethodPost, "/api/v1/signin", router.signIn)
}

type signInRequest struct {
	Username string `json:"username"`
	Passcode string `json:"passcode"`
}

func (s *sessionRouter) signIn(w http.ResponseWriter, r *http.Request) error {
	var request signInRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return fmt.Errorf("[handlers.signIn] failed to decode body with error: %v: %w", err, werrors.ErrBadRequest)
	}
	if request.Username == "" {
		return fmt.Errorf("[handlers.signIn] username is required: %w", werrors.ErrBadRequest)
	}
	res, err := s.controller.SignIn(request.Username, request.Passcode)
	if err != nil {
		return fmt.Errorf("[handlers.signIn] failed to sign in: %w", err)
	}
	web.Respond(r.Context(), w, map[string]any{"jwt": res}, http.StatusOK)
	return nil
}
