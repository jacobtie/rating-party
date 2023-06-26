package game

import "time"

type Game struct {
	GameID    string    `json:"gameId"`
	GameName  string    `json:"gameName"`
	GameCode  string    `json:"gameCode"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
