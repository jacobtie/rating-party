package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jacobtie/rating-party/server/internal/platform/config"
	"github.com/jacobtie/rating-party/server/internal/platform/contextvalue"
	"github.com/jacobtie/rating-party/server/internal/platform/web"
	"github.com/jacobtie/rating-party/server/internal/platform/werrors"
)

func AuthenticateMW(next web.Handler) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		cfg, err := config.Get()
		if err != nil {
			return fmt.Errorf("[middleware.AuthenticateMW] failed to get config: %w", err)
		}
		v, ok := r.Context().Value(contextvalue.KeyValues).(*contextvalue.Values)
		if !ok {
			return fmt.Errorf("[middleware.AuthenticateMW] failed to cast context values")
		}
		authHeader := r.Header.Get("Authorization")
		parsedToken, err := parseToken(cfg, authHeader)
		if err != nil {
			return err
		}
		if err := validateExpiration(parsedToken); err != nil {
			return err
		}
		v.JWT = parsedToken
		return next(w, r)
	}
}

func parseToken(cfg *config.Config, authHeader string) (*jwt.Token, error) {
	if authHeader == "" {
		return nil, fmt.Errorf("[middleware.AuthenticateMW] no auth header: %w", werrors.ErrUnauthorized)
	}
	authHeaderSegments := strings.Split(authHeader, " ")
	if len(authHeaderSegments) != 2 || authHeaderSegments[0] != "Bearer" {
		return nil, fmt.Errorf("[middleware.AuthenticateMW] malformed auth header: %w", werrors.ErrUnauthorized)
	}
	token := authHeaderSegments[1]
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method from JWT")
		}
		return []byte(cfg.AdminJWTSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("[middleware.AuthenticateMW] failed to parse token: %w", err)
	}
	return parsedToken, nil
}

func validateExpiration(parsedToken *jwt.Token) error {
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("[middleware.AuthenticateMW] failed to cast jwt claims: %w", werrors.ErrUnauthorized)
	}
	jwtEXPFieldRaw, ok := claims["exp"]
	if !ok {
		return fmt.Errorf("[middleware.AuthenticateMW] could not find exp on jwt: %w", werrors.ErrUnauthorized)
	}
	fmt.Printf("%+v\n", jwtEXPFieldRaw)
	jwtEXPField, ok := jwtEXPFieldRaw.(float64)
	if !ok {
		return fmt.Errorf("[middleware.AuthenticateMW] exp field was not a number: %w", werrors.ErrUnauthorized)
	}
	if int64(jwtEXPField) <= time.Now().Unix() {
		return fmt.Errorf("[middleware.AuthenticateMW] token is expired: %w", werrors.ErrUnauthorized)
	}
	return nil
}
