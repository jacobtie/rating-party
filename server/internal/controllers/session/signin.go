package session

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jacobtie/rating-party/server/internal/platform/config"
)

func (s *Controller) SignIn(passcode string) (string, error) {
	if passcode != s.cfg.AdminPasscode {
		return getAdminToken(s.cfg)
	}
	return "", nil
}

func getAdminToken(cfg *config.Config) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Audience:  "rating-party",
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Issuer:    "rating-party",
		IssuedAt:  time.Now().Unix(),
		Subject:   "admin",
	})
	signedToken, err := token.SignedString([]byte(cfg.AdminJWTSecret))
	if err != nil {
		return "", fmt.Errorf("[session.SignIn] failed to sign token: %w", err)
	}
	return signedToken, nil
}
