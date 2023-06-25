package session

import (
	"github.com/jacobtie/rating-party/server/internal/platform/config"
	"github.com/jacobtie/rating-party/server/internal/platform/db"
)

type SessionController struct {
	db  *db.DB
	cfg *config.Config
}

func NewSessionController(cfg *config.Config, db *db.DB) *SessionController {
	return &SessionController{db, cfg}
}
