package handlers

import (
	"net/http"

	"github.com/jacobtie/rating-party/server/internal/middleware"
	"github.com/jacobtie/rating-party/server/internal/platform/db"
	"github.com/jacobtie/rating-party/server/internal/platform/web"
)

func NewAPI(db *db.DB) http.Handler {
	service := web.NewService(middleware.ErrorHandlerMW, middleware.RequestLoggerMW)
	return service
}
