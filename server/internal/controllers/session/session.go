package session

import (
	"github.com/jacobtie/rating-party/server/internal/platform/config"
)

type Controller struct {
	cfg *config.Config
}

func NewController(cfg *config.Config) *Controller {
	return &Controller{cfg}
}
