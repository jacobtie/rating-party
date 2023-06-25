package middleware

import (
	"net/http"

	"github.com/jacobtie/rating-party/server/internal/platform/contextvalue"
	"github.com/jacobtie/rating-party/server/internal/platform/logger"
	"github.com/jacobtie/rating-party/server/internal/platform/web"
)

// requestLoggerMW acts as a "built-in" middleware to log requests
func RequestLoggerMW(before web.Handler) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := before(w, r)
		ctx := r.Context()
		v, ok := ctx.Value(contextvalue.KeyValues).(*contextvalue.Values)
		if !ok {
			return err
		}
		logger.GetCtx(r.Context()).Debug().
			Int("status", v.StatusCode).
			Bool("error", v.Error).
			Msg("request log")
		return err
	}
}
