package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/jacobtie/rating-party/server/internal/platform/contextvalue"
	"github.com/jacobtie/rating-party/server/internal/platform/web"
	"github.com/jacobtie/rating-party/server/internal/platform/werrors"
)

func MakeAuthorizationMW(scopes ...string) web.Middleware {
	return func(next web.Handler) web.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			v, ok := r.Context().Value(contextvalue.KeyValues).(*contextvalue.Values)
			if !ok {
				return fmt.Errorf("[middleware.MakeAuthorizationMW] failed to cast context values")
			}
			token := v.JWT
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return fmt.Errorf("[middleware.MakeAuthorizationMW] failed to cast jwt claims: %w", werrors.ErrForbidden)
			}
			isAdmin, err := isJWTAdmin(claims)
			if err != nil {
				return err
			}
			if isAdmin {
				return next(w, r)
			}
			if err := checkScopes(scopes, claims); err != nil {
				return err
			}
			return next(w, r)
		}
	}
}

func isJWTAdmin(claims jwt.MapClaims) (bool, error) {
	jwtSubFieldRaw, ok := claims["sub"]
	if !ok {
		return false, nil
	}
	jwtSubField, ok := jwtSubFieldRaw.(string)
	if !ok {
		return false, fmt.Errorf("[middleware.MakeAuthorizationMW] failed to parse sub field as a string: %w", werrors.ErrForbidden)
	}
	if jwtSubField != "admin" {
		return false, nil
	}
	return true, nil
}

func checkScopes(scopes []string, claims jwt.MapClaims) error {
	jwtScopesFieldRaw, ok := claims["scope"]
	if !ok {
		return fmt.Errorf("[middleware.MakeAuthorizationMW] could not find scope on jwt: %w", werrors.ErrForbidden)
	}
	jwtScopesField, ok := jwtScopesFieldRaw.(string)
	if !ok {
		return fmt.Errorf("[middleware.MakeAuthorizationMW] scope field was not a string: %w", werrors.ErrForbidden)
	}
	jwtScopes := strings.Split(jwtScopesField, " ")
	scopesSet := make(map[string]struct{}, len(jwtScopes))
	for _, scope := range jwtScopes {
		scopesSet[scope] = struct{}{}
	}
	for _, scope := range scopes {
		if _, ok := scopesSet[scope]; !ok {
			return fmt.Errorf("[middleware.MakeAuthorizationMW] missing required scope: %w", werrors.ErrForbidden)
		}
	}
	return nil
}
