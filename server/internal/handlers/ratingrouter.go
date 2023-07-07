package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jacobtie/rating-party/server/internal/config"
	"github.com/jacobtie/rating-party/server/internal/controllers/rating"
	"github.com/jacobtie/rating-party/server/internal/middleware"
	"github.com/jacobtie/rating-party/server/internal/platform/contextvalue"
	"github.com/jacobtie/rating-party/server/internal/platform/db"
	"github.com/jacobtie/rating-party/server/internal/platform/web"
	"github.com/jacobtie/rating-party/server/internal/platform/werrors"
	"github.com/julienschmidt/httprouter"
)

type ratingRouter struct {
	controller *rating.Controller
}

func registerRatingRoutes(service *web.Service, cfg *config.Config, db *db.DB) {
	router := &ratingRouter{
		controller: rating.NewController(cfg, db),
	}
	service.Handle(http.MethodGet, "/api/v1/games/:gameId/ratings", router.getRatings, middleware.MakeAuthorizationMW(false), middleware.AuthenticateMW)
	service.Handle(http.MethodPut, "/api/v1/games/:gameId/wines/:wineId/ratings", router.putRating, middleware.MakeAuthorizationMW(false), middleware.AuthenticateMW)
}

func (rr *ratingRouter) getRatings(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	if params == nil {
		return fmt.Errorf("[handlers.getRating] no params in context: %w", werrors.ErrBadRequest)
	}
	gameID := params.ByName("gameId")
	if gameID == "" {
		return fmt.Errorf("[handlers.getRating] game ID was not found: %w", werrors.ErrBadRequest)
	}
	if _, err := uuid.Parse(gameID); err != nil {
		return fmt.Errorf("[handlers.getRating] game ID was not a UUID: %w", werrors.ErrNotFound)
	}
	v, ok := ctx.Value(contextvalue.KeyValues).(*contextvalue.Values)
	if !ok {
		return fmt.Errorf("[handlers.getRating] no values in context")
	}
	if v.IsAdmin {
		ratings, err := rr.controller.GetAllByGameID(ctx, gameID)
		if err != nil {
			return fmt.Errorf("[handlers.getRating]: could not get all ratings for admin: %w", err)
		}
		web.Respond(ctx, w, ratings, http.StatusOK)
		return nil
	}
	if v.UserID == "" {
		return fmt.Errorf("[handlers.getRating] user ID was not found")
	}
	ratings, err := rr.controller.GetAllByGameIDAndParticipantID(ctx, gameID, v.UserID)
	if err != nil {
		return fmt.Errorf("[handlers.getRating]: could not get all ratings for user: %w", err)
	}
	web.Respond(ctx, w, ratings, http.StatusOK)
	return nil
}

type putRatingRequest struct {
	SightRating   int    `json:"sightRating"`
	AromaRating   int    `json:"aromaRating"`
	TasteRating   int    `json:"tasteRating"`
	OverallRating int    `json:"overallRating"`
	Comments      string `json:"comments"`
}

func (rr *ratingRouter) putRating(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	if params == nil {
		return fmt.Errorf("[handlers.putRating] no params in context: %w", werrors.ErrBadRequest)
	}
	gameID := params.ByName("gameId")
	if gameID == "" {
		return fmt.Errorf("[handlers.putRating] game ID was not found: %w", werrors.ErrBadRequest)
	}
	if _, err := uuid.Parse(gameID); err != nil {
		return fmt.Errorf("[handlers.putRating] game ID was not a UUID: %w", werrors.ErrNotFound)
	}
	wineID := params.ByName("wineId")
	if wineID == "" {
		return fmt.Errorf("[handlers.putRating] wine ID was not found: %w", werrors.ErrBadRequest)
	}
	if _, err := uuid.Parse(wineID); err != nil {
		return fmt.Errorf("[handlers.putRating] wine ID was not a UUID: %w", werrors.ErrNotFound)
	}
	v, ok := ctx.Value(contextvalue.KeyValues).(*contextvalue.Values)
	if !ok {
		return fmt.Errorf("[handlers.putRating] no values in context")
	}
	if v.UserID == "" {
		return fmt.Errorf("[handlers.putRating] user ID was not found")
	}
	var req putRatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("[handlers.putRating] failed to decode request body: %w", werrors.ErrBadRequest)
	}
	rating, err := rr.controller.UpsertRating(ctx, &rating.Rating{
		ParticipantID: v.UserID,
		GameID:        gameID,
		WineID:        wineID,
		SightRating:   req.SightRating,
		AromaRating:   req.AromaRating,
		TasteRating:   req.TasteRating,
		OverallRating: req.OverallRating,
		Comments:      req.Comments,
	})
	if err != nil {
		return fmt.Errorf("[handlers.putRating]: %w", err)
	}
	web.Respond(ctx, w, rating, http.StatusOK)
	return nil
}
