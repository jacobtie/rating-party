package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jacobtie/rating-party/server/internal/config"
	"github.com/jacobtie/rating-party/server/internal/controllers/game"
	"github.com/jacobtie/rating-party/server/internal/middleware"
	"github.com/jacobtie/rating-party/server/internal/platform/db"
	"github.com/jacobtie/rating-party/server/internal/platform/web"
	"github.com/jacobtie/rating-party/server/internal/platform/werrors"
	"github.com/julienschmidt/httprouter"
)

type gameRouter struct {
	controller *game.Controller
}

func registerGameRoutes(service *web.Service, cfg *config.Config, db *db.DB) {
	router := &gameRouter{
		controller: game.NewController(cfg, db),
	}
	service.Handle(http.MethodGet, "/api/v1/games", router.getAllGames, middleware.MakeAuthorizationMW(true), middleware.AuthenticateMW)
	service.Handle(http.MethodGet, "/api/v1/games/:gameId", router.getSingleGame, middleware.MakeAuthorizationMW(false), middleware.AuthenticateMW)
	service.Handle(http.MethodPost, "/api/v1/games", router.createGame, middleware.MakeAuthorizationMW(true), middleware.AuthenticateMW)
	service.Handle(http.MethodPut, "/api/v1/games/:gameId", router.updateGame, middleware.MakeAuthorizationMW(true), middleware.AuthenticateMW)
	service.Handle(http.MethodDelete, "/api/v1/games/:gameId", router.deleteGame, middleware.MakeAuthorizationMW(true), middleware.AuthenticateMW)
}

func (g *gameRouter) getAllGames(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	games, err := g.controller.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("[handlers.getAllGames]: %w", err)
	}
	web.Respond(ctx, w, games, http.StatusOK)
	return nil
}

func (g *gameRouter) getSingleGame(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	if params == nil {
		return fmt.Errorf("[handlers.getSingleGame] no params in context: %w", werrors.ErrBadRequest)
	}
	gameID := params.ByName("gameId")
	if gameID == "" {
		return fmt.Errorf("[handlers.getSingleGame] game ID was not found: %w", werrors.ErrBadRequest)
	}
	if _, err := uuid.Parse(gameID); err != nil {
		return fmt.Errorf("[handlers.getSingleGame] game ID was not a UUID: %w", werrors.ErrBadRequest)
	}
	game, err := g.controller.GetSingle(ctx, gameID)
	if err != nil {
		return fmt.Errorf("[handlers.getSingleGame]: %w", err)
	}
	web.Respond(ctx, w, game, http.StatusOK)
	return nil
}

type createGameRequest struct {
	GameName string `json:"gameName"`
}

func (g *gameRouter) createGame(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req createGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("[handlers.createGame] failed to decode request: %w", werrors.ErrBadRequest)
	}
	if req.GameName == "" {
		return fmt.Errorf("[handlers.createGame] game name was empty: %w", werrors.ErrBadRequest)
	}
	game, err := g.controller.Create(ctx, req.GameName)
	if err != nil {
		return fmt.Errorf("[handlers.createGame]: %w", err)
	}
	web.Respond(ctx, w, game, http.StatusCreated)
	return nil
}

type updateGameRequest struct {
	GameName  string `json:"gameName"`
	IsRunning bool   `json:"isRunning"`
}

func (g *gameRouter) updateGame(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	if params == nil {
		return fmt.Errorf("[handlers.updateGame] no params in context: %w", werrors.ErrBadRequest)
	}
	gameID := params.ByName("gameId")
	if gameID == "" {
		return fmt.Errorf("[handlers.updateGame] game ID was not found: %w", werrors.ErrBadRequest)
	}
	if _, err := uuid.Parse(gameID); err != nil {
		return fmt.Errorf("[handlers.updateGame] game ID was not a UUID: %w", werrors.ErrBadRequest)
	}
	var req updateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("[handlers.updateGame] failed to decode request: %w", werrors.ErrBadRequest)
	}
	if req.GameName == "" {
		return fmt.Errorf("[handlers.updateGame] game name was empty: %w", werrors.ErrBadRequest)
	}
	if err := g.controller.Update(ctx, gameID, req.GameName, req.IsRunning); err != nil {
		return fmt.Errorf("[handlers.updateGame]: %w", err)
	}
	web.Respond(ctx, w, nil, http.StatusNoContent)
	return nil
}

func (g *gameRouter) deleteGame(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	if params == nil {
		return fmt.Errorf("[handlers.deleteGame] no params in context: %w", werrors.ErrBadRequest)
	}
	gameID := params.ByName("gameId")
	if gameID == "" {
		return fmt.Errorf("[handlers.deleteGame] game ID was not found: %w", werrors.ErrBadRequest)
	}
	if _, err := uuid.Parse(gameID); err != nil {
		return fmt.Errorf("[handlers.deleteGame] game ID was not a UUID: %w", werrors.ErrBadRequest)
	}
	game, err := g.controller.Delete(ctx, gameID)
	if err != nil {
		return fmt.Errorf("[handlers.deleteGame]: %w", err)
	}
	web.Respond(ctx, w, game, http.StatusOK)
	return nil
}
