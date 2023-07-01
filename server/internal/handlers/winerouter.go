package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jacobtie/rating-party/server/internal/config"
	"github.com/jacobtie/rating-party/server/internal/controllers/wine"
	"github.com/jacobtie/rating-party/server/internal/middleware"
	"github.com/jacobtie/rating-party/server/internal/platform/db"
	"github.com/jacobtie/rating-party/server/internal/platform/web"
	"github.com/jacobtie/rating-party/server/internal/platform/werrors"
	"github.com/julienschmidt/httprouter"
)

type wineRouter struct {
	controller *wine.Controller
}

func registerWineRoutes(service *web.Service, cfg *config.Config, db *db.DB) {
	router := &wineRouter{
		controller: wine.NewController(cfg, db),
	}
	service.Handle(http.MethodGet, "/api/v1/games/:gameId/wines", router.getAllWines, middleware.MakeAuthorizationMW(false), middleware.AuthenticateMW)
	service.Handle(http.MethodGet, "/api/v1/games/:gameId/wines/:wineId", router.getSingleWine, middleware.MakeAuthorizationMW(false), middleware.AuthenticateMW)
	service.Handle(http.MethodPost, "/api/v1/games/:gameId/wines", router.createWine, middleware.MakeAuthorizationMW(true), middleware.AuthenticateMW)
	service.Handle(http.MethodPut, "/api/v1/games/:gameId/wines/:wineId", router.updateWine, middleware.MakeAuthorizationMW(true), middleware.AuthenticateMW)
	service.Handle(http.MethodDelete, "/api/v1/games/:gameId/wines/:wineId", router.deleteWine, middleware.MakeAuthorizationMW(true), middleware.AuthenticateMW)
}

func (wr *wineRouter) getAllWines(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	gameID := params.ByName("gameId")
	wines, err := wr.controller.GetAllWines(ctx, gameID)
	if err != nil {
		return fmt.Errorf("[wineRouter.getAllWines] failed to get all wines: %w", err)
	}
	web.Respond(ctx, w, wines, http.StatusOK)
	return nil
}

func (wr *wineRouter) getSingleWine(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	wineID := params.ByName("wineId")
	wine, err := wr.controller.GetSingleWine(ctx, wineID)
	if err != nil {
		return fmt.Errorf("[wineRouter.getSingleWine] failed to get single wine: %w", err)
	}
	web.Respond(ctx, w, wine, http.StatusOK)
	return nil
}

type createUpdateWineRequest struct {
	WineName string `json:"wineName"`
	WineCode string `json:"wineCode"`
	WineYear int    `json:"wineYear"`
}

func (wr *wineRouter) createWine(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	gameID := params.ByName("gameId")
	var req createUpdateWineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("[wineRouter.createWine] failed to decode request body: %w", werrors.ErrBadRequest)
	}
	if req.WineName == "" {
		return fmt.Errorf("[wineRouter.createWine] wine name is required: %w", werrors.ErrBadRequest)
	}
	if req.WineCode == "" {
		return fmt.Errorf("[wineRouter.createWine] wine code is required: %w", werrors.ErrBadRequest)
	}
	if req.WineYear == 0 {
		return fmt.Errorf("[wineRouter.createWine] wine year is required: %w", werrors.ErrBadRequest)
	}
	wine, err := wr.controller.CreateWine(ctx, gameID, req.WineName, req.WineCode, req.WineYear)
	if err != nil {
		return fmt.Errorf("[wineRouter.createWine] failed to create wine: %w", err)
	}
	web.Respond(ctx, w, wine, http.StatusCreated)
	return nil
}

func (wr *wineRouter) updateWine(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	wineID := params.ByName("wineId")
	var req createUpdateWineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("[wineRouter.updateWine] failed to decode request body: %w", werrors.ErrBadRequest)
	}
	if req.WineName == "" {
		return fmt.Errorf("[wineRouter.updateWine] wine name is required: %w", werrors.ErrBadRequest)
	}
	if req.WineCode == "" {
		return fmt.Errorf("[wineRouter.updateWine] wine code is required: %w", werrors.ErrBadRequest)
	}
	if req.WineYear == 0 {
		return fmt.Errorf("[wineRouter.updateWine] wine year is required: %w", werrors.ErrBadRequest)
	}
	if err := wr.controller.UpdateWine(ctx, wineID, req.WineName, req.WineCode, req.WineYear); err != nil {
		return fmt.Errorf("[wineRouter.updateWine] failed to update wine: %w", err)
	}
	web.Respond(ctx, w, nil, http.StatusNoContent)
	return nil
}

func (wr *wineRouter) deleteWine(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	wineID := params.ByName("wineId")
	wine, err := wr.controller.DeleteWine(ctx, wineID)
	if err != nil {
		return fmt.Errorf("[wineRouter.deleteWine] failed to delete wine: %w", err)
	}
	web.Respond(ctx, w, wine, http.StatusOK)
	return nil
}
