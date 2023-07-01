package session

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func (s *Controller) SignIn(username, passcode string) (string, error) {
	if username == "admin" && passcode == s.cfg.AdminPasscode {
		return s.signAdminToken()
	}
	return s.signInToGame(username, passcode)
}

func (s *Controller) signAdminToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Audience:  "rating-party",
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Issuer:    "rating-party",
		IssuedAt:  time.Now().Unix(),
		Subject:   "admin",
	})
	signedToken, err := token.SignedString([]byte(s.cfg.AdminJWTSecret))
	if err != nil {
		return "", fmt.Errorf("[session.signAdminToken] failed to sign token: %w", err)
	}
	return signedToken, nil
}
