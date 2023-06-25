package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jacobtie/rating-party/server/internal/controllers/session"
	"github.com/jacobtie/rating-party/server/internal/platform/config"
	"github.com/jacobtie/rating-party/server/internal/platform/db"
	"github.com/jacobtie/rating-party/server/internal/platform/web"
	"github.com/jacobtie/rating-party/server/internal/platform/werrors.go"
)

type sessionRouter struct {
	sessionController *session.SessionController
}

func registerSessionRoutes(server *web.Service, cfg *config.Config, db *db.DB) {
	sessionRouter := &sessionRouter{
		sessionController: session.NewSessionController(cfg, db),
	}
	server.Handle(http.MethodPost, "/signin", sessionRouter.SignIn)
}

type SignInRequest struct {
	Passcode string `json:"passcode"`
}

func (s *sessionRouter) SignIn(w http.ResponseWriter, r *http.Request) error {
	var request SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return fmt.Errorf("[session.SignIn] failed to decode body with error: %v: %w", err, werrors.ErrBadRequest)
	}
	jwt, err := s.sessionController.SignIn(request.Passcode)
	if err != nil {
		return fmt.Errorf("[session.SignIn] failed to sign in: %w", err)
	}
	web.Respond(r.Context(), w, map[string]any{"jwt": jwt}, http.StatusOK)
	return nil
}
