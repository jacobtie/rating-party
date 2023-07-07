package session

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type SignInResponse struct {
	JWT     string  `json:"jwt"`
	IsAdmin *bool   `json:"isAdmin,omitempty"`
	GameID  *string `json:"gameId,omitempty"`
}

func (s *Controller) SignIn(ctx context.Context, username, passcode string) (*SignInResponse, error) {
	if username == "admin" && passcode == s.cfg.AdminPasscode {
		return s.signAdminToken()
	}
	return s.signInToGame(ctx, username, passcode)
}

func (s *Controller) signAdminToken() (*SignInResponse, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Audience:  "rating-party",
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Issuer:    "rating-party",
		IssuedAt:  time.Now().Unix(),
		Subject:   "admin",
	})
	signedToken, err := token.SignedString([]byte(s.cfg.AdminJWTSecret))
	if err != nil {
		return nil, fmt.Errorf("[session.signAdminToken] failed to sign token: %w", err)
	}
	isAdmin := true
	return &SignInResponse{
		JWT:     signedToken,
		IsAdmin: &isAdmin,
	}, nil
}
