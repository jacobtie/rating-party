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

type CreateUpdateGameRequest struct {
	GameName string `json:"gameName"`
}

func (g *gameRouter) CreateGame(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req CreateUpdateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("[handlers.CreateGame] failed to decode request: %w", werrors.ErrBadRequest)
	}
	if req.GameName == "" {
		return fmt.Errorf("[handlers.CreateGame] game name was empty: %w", werrors.ErrBadRequest)
	}
	game, err := g.controller.Create(ctx, req.GameName)
	if err != nil {
		return fmt.Errorf("[handlers.CreateGame]: %w", err)
	}
	web.Respond(ctx, w, game, http.StatusOK)
	return nil
}

func (g *gameRouter) UpdateGame(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	if params == nil {
		return fmt.Errorf("[handlers.UpdateGame] no params in context: %w", werrors.ErrBadRequest)
	}
	gameID := params.ByName("gameId")
	if gameID == "" {
		return fmt.Errorf("[handlers.UpdateGame] game ID was not found: %w", werrors.ErrBadRequest)
	}
	if _, err := uuid.Parse(gameID); err != nil {
		return fmt.Errorf("[handlers.UpdateGame] game ID was not a UUID: %w", werrors.ErrNotFound)
	}
	var req CreateUpdateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("[handlers.UpdateGame] failed to decode request: %w", werrors.ErrBadRequest)
	}
	if req.GameName == "" {
		return fmt.Errorf("[handlers.UpdateGame] game name was empty: %w", werrors.ErrBadRequest)
	}
	if err := g.controller.Update(ctx, gameID, req.GameName); err != nil {
		return fmt.Errorf("[handlers.UpdateGame]: %w", err)
	}
	web.Respond(ctx, w, nil, http.StatusNoContent)
	return nil
}

func (g *gameRouter) DeleteGame(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	if params == nil {
		return fmt.Errorf("[handlers.DeleteGame] no params in context: %w", werrors.ErrBadRequest)
	}
	gameID := params.ByName("gameId")
	if gameID == "" {
		return fmt.Errorf("[handlers.DeleteGame] game ID was not found: %w", werrors.ErrBadRequest)
	}
	if _, err := uuid.Parse(gameID); err != nil {
		return fmt.Errorf("[handlers.DeleteGame] game ID was not a UUID: %w", werrors.ErrNotFound)
	}
	game, err := g.controller.Delete(ctx, gameID)
	if err != nil {
		return fmt.Errorf("[handlers.DeleteGame]: %w", err)
	}
	web.Respond(ctx, w, game, http.StatusOK)
	return nil
}
