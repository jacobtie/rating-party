package session

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jacobtie/rating-party/server/internal/platform/werrors"
)

func (s *Controller) signInToGame(username, passcode string) (string, error) {
	gameID, err := s.getGameID(passcode)
	if err != nil {
		return "", fmt.Errorf("[session.signInToGame] failed to get game id: %w", err)
	}
	token, err := s.signUserToken(username, gameID)
	if err != nil {
		return "", fmt.Errorf("[session.signInToGame] failed to sign token: %w", err)
	}
	return token, nil
}

func (s *Controller) getGameID(passcode string) (string, error) {
	row := s.db.QueryRow("SELECT BIN_TO_UUID(game_id) FROM game WHERE game_code = ?", passcode)
	var gameID string
	if err := row.Scan(&gameID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("[session.getGameID] no game found: %w", werrors.ErrUnauthorized)
		}
		return "", err
	}
	return gameID, nil
}

func (s *Controller) signUserToken(username, gameID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud":    "rating-party",
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"iss":    "rating-party",
		"iat":    time.Now().Unix(),
		"sub":    username,
		"gameId": gameID,
	})
	signedToken, err := token.SignedString([]byte(s.cfg.AdminJWTSecret))
	if err != nil {
		return "", fmt.Errorf("[session.signUserToken] failed to sign token: %w", err)
	}
	return signedToken, nil
}
