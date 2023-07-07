package rating

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/jacobtie/rating-party/server/internal/config"
	"github.com/jacobtie/rating-party/server/internal/controllers/participant"
	"github.com/jacobtie/rating-party/server/internal/controllers/wine"
	"github.com/jacobtie/rating-party/server/internal/platform/db"
)

type Controller struct {
	cfg                   *config.Config
	db                    *db.DB
	participantController *participant.Controller
	wineController        *wine.Controller
}

func NewController(cfg *config.Config, db *db.DB, participantController *participant.Controller, wineController *wine.Controller) *Controller {
	return &Controller{
		cfg:                   cfg,
		db:                    db,
		participantController: participantController,
		wineController:        wineController,
	}
}

func (c *Controller) GetRatingsResult(ctx context.Context, gameID string, includeUsernames bool) ([]map[string]any, error) {
	ratings, err := c.GetAllByGameID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("[controllers.rating.GetRatingsResult] could not get all ratings: %w", err)
	}
	participants, err := c.participantController.GetAllParticipantsByGameID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("[controllers.rating.GetRatingsResult] could not get all participants: %w", err)
	}
	usernames := make([]string, len(participants))
	for _, participant := range participants {
		usernames = append(usernames, participant.Username)
	}
	wineAggMap := make(map[string]map[string]float64)
	for _, rating := range ratings {
		// Skip empty ratings
		if rating.SightRating == 0 &&
			rating.AromaRating == 0 &&
			rating.TasteRating == 0 &&
			rating.OverallRating == 0 &&
			rating.Comments == "" {
			continue
		}
		if _, ok := wineAggMap[rating.WineID]; !ok {
			wineAggMap[rating.WineID] = make(map[string]float64)
		}
		wineAggMap[rating.WineID][rating.Username] = rating.SightRating + rating.AromaRating + rating.TasteRating + rating.OverallRating
	}
	rows := make([]map[string]any, 0)
	for wineID, wineScores := range wineAggMap {
		wine, err := c.wineController.GetSingleWine(ctx, wineID)
		if err != nil {
			return nil, fmt.Errorf("[controllers.rating.GetRatingsResult] could not get wine by id: %w", err)
		}
		row := make(map[string]any)
		row["wineID"] = wine.WineID
		row["wineName"] = wine.WineName
		row["wineCode"] = wine.WineCode
		row["wineYear"] = wine.WineYear
		sum := 0.0
		count := 0
		for _, username := range usernames {
			score, ok := wineScores[username]
			if !ok {
				score = 0.0
			}
			if includeUsernames {
				row[username] = score
			}
			if !ok {
				continue
			}
			sum += score
			count += 1
		}
		if count == 0 {
			row["avg"] = 0.0
		} else {
			row["avg"] = math.Round((sum/float64(count))*100) / 100
		}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i]["avg"].(float64) > rows[j]["avg"].(float64)
	})
	for i := 0; i < len(rows); i++ {
		rows[i]["rank"] = i + 1
	}
	return rows, nil
}
