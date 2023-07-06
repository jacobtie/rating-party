package handlers

import (
	"fmt"
	"net/http"

	"github.com/jacobtie/rating-party/server/internal/config"
	"github.com/jacobtie/rating-party/server/internal/controllers/participant"
	"github.com/jacobtie/rating-party/server/internal/middleware"
	"github.com/jacobtie/rating-party/server/internal/platform/db"
	"github.com/jacobtie/rating-party/server/internal/platform/web"
	"github.com/jacobtie/rating-party/server/internal/platform/werrors"
	"github.com/julienschmidt/httprouter"
)

type participantRouter struct {
	controller *participant.Controller
}

func registerParticipantRoutes(service *web.Service, cfg *config.Config, db *db.DB) {
	router := participantRouter{
		controller: participant.NewController(cfg, db),
	}
	service.Handle(http.MethodGet, "/api/v1/games/:gameId/participants", router.getAllParticipants, middleware.MakeAuthorizationMW(true), middleware.AuthenticateMW)
}

func (p *participantRouter) getAllParticipants(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	if params == nil {
		return fmt.Errorf("[handlers.getAllParticipants] no params in context: %w", werrors.ErrBadRequest)
	}
	gameID := params.ByName("gameId")
	if gameID == "" {
		return fmt.Errorf("[handlers.getAllParticipants] game ID was not found: %w", werrors.ErrBadRequest)
	}
	participants, err := p.controller.GetAllParticipantsByGameID(ctx, gameID)
	if err != nil {
		return fmt.Errorf("[handlers.getAllParticipants]: %w", err)
	}
	web.Respond(ctx, w, participants, http.StatusOK)
	return nil
}
