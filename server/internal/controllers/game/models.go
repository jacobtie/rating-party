package game

type Game struct {
	GameID    string `json:"gameId"`
	GameName  string `json:"gameName"`
	GameCode  string `json:"gameCode"`
	IsRunning bool   `json:"isRunning"`
}
