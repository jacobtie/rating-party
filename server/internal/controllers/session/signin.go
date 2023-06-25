package session

import (
	"fmt"
	"time"

	"github.com/jacobtie/rating-party/server/internal/platform/werrors.go"

	"github.com/golang-jwt/jwt"
)

func (s *SessionController) SignIn(passcode string) (string, error) {
	if passcode != s.cfg.AdminPasscode {
		return "", fmt.Errorf("[session.SignIn] invalid passcode: %w", werrors.ErrUnauthorized)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Audience:  "ranking-party",
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Issuer:    "ranking-party",
		IssuedAt:  time.Now().Unix(),
		Subject:   "admin",
	})
	signedToken, err := token.SignedString([]byte(s.cfg.AdminJWTSecret))
	if err != nil {
		return "", fmt.Errorf("[session.SignIn] failed to sign token: %w", err)
	}
	return signedToken, nil
}
