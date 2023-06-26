package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jacobtie/rating-party/server/internal/controllers/game"
	"github.com/jacobtie/rating-party/server/internal/middleware"
	"github.com/jacobtie/rating-party/server/internal/platform/config"
	"github.com/jacobtie/rating-party/server/internal/platform/db"
	"github.com/jacobtie/rating-party/server/internal/platform/web"
	"github.com/jacobtie/rating-party/server/internal/platform/werrors"
	"github.com/julienschmidt/httprouter"
)

type gameRouter struct {
	controller *game.Controller
}

func registerGameRoutes(service *web.Service, db *db.DB, cfg *config.Config) {
	router := &gameRouter{
		controller: game.NewController(cfg, db),
	}
	service.Handle(http.MethodGet, "/api/v1/games", router.GetAllGames, middleware.MakeAuthorizationMW("read:games"), middleware.AuthenticateMW)
	service.Handle(http.MethodGet, "/api/v1/games/:gameId", router.GetSingleGame, middleware.MakeAuthorizationMW("read:game"), middleware.AuthenticateMW)
	service.Handle(http.MethodPost, "/api/v1/games", router.CreateGame, middleware.MakeAuthorizationMW("create:game"), middleware.AuthenticateMW)
	service.Handle(http.MethodPut, "/api/v1/games/:gameId", router.UpdateGame, middleware.MakeAuthorizationMW("update:game"), middleware.AuthenticateMW)
	service.Handle(http.MethodDelete, "/api/v1/games/:gameId", router.DeleteGame, middleware.MakeAuthorizationMW("delete:game"), middleware.AuthenticateMW)
}

func (g *gameRouter) GetAllGames(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	games, err := g.controller.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("[handlers.GetAllGames]: %w", err)
	}
	web.Respond(ctx, w, games, http.StatusOK)
	return nil
}

func (g *gameRouter) GetSingleGame(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	if params == nil {
		return fmt.Errorf("[handlers.GetSingleGame] no params in context: %w", werrors.ErrBadRequest)
	}
	gameID := params.ByName("gameId")
	if gameID == "" {
		return fmt.Errorf("[handlers.GetSingleGame] game ID was not found: %w", werrors.ErrBadRequest)
	}
	if _, err := uuid.Parse(gameID); err != nil {
		return fmt.Errorf("[handlers.GetSingleGame] game ID was not a UUID: %w", werrors.ErrNotFound)
	}
	game, err := g.controller.GetSingle(ctx, gameID)
	if err != nil {
		return fmt.Errorf("[handlers.GetSingleGame]: %w", err)
	}
	web.Respond(ctx, w, game, http.StatusOK)
	return nil
}

func (g *gameRouter) CreateGame(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (g *gameRouter) UpdateGame(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (g *gameRouter) DeleteGame(w http.ResponseWriter, r *http.Request) error {
	return nil
}
