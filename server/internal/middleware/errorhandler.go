package middleware

import (
	"net/http"

	"github.com/jacobtie/rating-party/server/internal/platform/web"
)

func ErrorHandlerMW(before web.Handler) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := before(w, r); err != nil {
			web.HandleError(r.Context(), w, err)
		}
		return nil
	}
}
