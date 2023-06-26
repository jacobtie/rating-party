package handlers

import (
	"net/http"

	"github.com/jacobtie/rating-party/server/internal/middleware"
	"github.com/jacobtie/rating-party/server/internal/platform/config"
	"github.com/jacobtie/rating-party/server/internal/platform/db"
	"github.com/jacobtie/rating-party/server/internal/platform/web"
)

func NewAPI(cfg *config.Config, db *db.DB) http.Handler {
	service := web.NewService(middleware.ErrorHandlerMW, middleware.RequestLoggerMW)
	registerSessionRoutes(service, cfg)
	registerGameRoutes(service, db, cfg)
	return service
}
