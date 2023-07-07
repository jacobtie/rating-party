package rating

import (
	"github.com/jacobtie/rating-party/server/internal/config"
	"github.com/jacobtie/rating-party/server/internal/platform/db"
)

type Controller struct {
	cfg *config.Config
	db  *db.DB
}

func NewController(cfg *config.Config, db *db.DB) *Controller {
	return &Controller{
		cfg: cfg,
		db:  db,
	}
}
