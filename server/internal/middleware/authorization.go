package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/jacobtie/rating-party/server/internal/platform/contextvalue"
	"github.com/jacobtie/rating-party/server/internal/platform/web"
	"github.com/jacobtie/rating-party/server/internal/platform/werrors"
	"github.com/julienschmidt/httprouter"
)

func MakeAuthorizationMW(requiresAdmin bool) web.Middleware {
	return func(next web.Handler) web.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()
			v, ok := ctx.Value(contextvalue.KeyValues).(*contextvalue.Values)
			if !ok {
				return fmt.Errorf("[middleware.MakeAuthorizationMW] failed to cast context values")
			}
			token := v.JWT
			if token == nil {
				return fmt.Errorf("[middleware.MakeAuthorizationMW] no JWT in context, must call AuthenticateMW first")
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return fmt.Errorf("[middleware.MakeAuthorizationMW] failed to cast jwt claims: %w", werrors.ErrForbidden)
			}
			isAdmin, userID, err := getIsJWTAdminOrUserID(claims)
			if err != nil {
				return err
			}
			if isAdmin {
				v.IsAdmin = true
				return next(w, r)
			}
			if !isAdmin && requiresAdmin {
				return fmt.Errorf("[middleware.MakeAuthorizationMW] requires admin: %w", werrors.ErrForbidden)
			}
			params := httprouter.ParamsFromContext(ctx)
			if params == nil {
				return fmt.Errorf("[middleware.MakeAuthorizationMW] failed to get params from context: %w", werrors.ErrForbidden)
			}
			gameID := params.ByName("gameId")
			if err := checkScopes(gameID, claims); err != nil {
				return err
			}
			v.UserID = userID
			return next(w, r)
		}
	}
}

func getIsJWTAdminOrUserID(claims jwt.MapClaims) (bool, string, error) {
	jwtSubFieldRaw, ok := claims["sub"]
	if !ok {
		return false, "", nil
	}
	jwtSubField, ok := jwtSubFieldRaw.(string)
	if !ok {
		return false, "", fmt.Errorf("[middleware.MakeAuthorizationMW] failed to parse sub field as a string: %w", werrors.ErrForbidden)
	}
	if jwtSubField != "admin" {
		return false, jwtSubField, nil
	}
	return true, "", nil
}

func checkScopes(gameID string, claims jwt.MapClaims) error {
	gameIDFieldRaw, ok := claims["gameId"]
	if !ok {
		return fmt.Errorf("[middleware.MakeAuthorizationMW] could not find gameId on jwt: %w", werrors.ErrForbidden)
	}
	gameIDField, ok := gameIDFieldRaw.(string)
	if !ok {
		return fmt.Errorf("[middleware.MakeAuthorizationMW] gameId field was not a string: %w", werrors.ErrForbidden)
	}
	if gameIDField != gameID {
		return fmt.Errorf("[middleware.MakeAuthorizationMW] gameId on jwt does not match gameId in url: %w", werrors.ErrForbidden)
	}
	return nil
}
