package game

import (
	"github.com/jacobtie/rating-party/server/internal/platform/config"
	"github.com/jacobtie/rating-party/server/internal/platform/db"
)

type Controller struct {
	db  *db.DB
	cfg *config.Config
}

func NewController(cfg *config.Config, db *db.DB) *Controller {
	return &Controller{db, cfg}
}
