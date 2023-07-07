package handlers

import (
	"io/fs"
	"net/http"

	"github.com/jacobtie/rating-party/server/internal/config"
	"github.com/jacobtie/rating-party/server/internal/middleware"
	"github.com/jacobtie/rating-party/server/internal/platform/db"
	"github.com/jacobtie/rating-party/server/internal/platform/web"
)

func NewAPI(cfg *config.Config, db *db.DB, clientDir fs.FS) http.Handler {
	service := web.NewService(middleware.ErrorHandlerMW, middleware.RequestLoggerMW)
	registerSessionRoutes(service, cfg, db)
	registerGameRoutes(service, cfg, db)
	registerWineRoutes(service, cfg, db)
	registerParticipantRoutes(service, cfg, db)
	registerRatingRoutes(service, cfg, db)
	// Serve SPA
	service.ServeSPA(clientDir)
	return service
}
