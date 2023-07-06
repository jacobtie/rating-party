package participant

type Participant struct {
	ParticipantID string `json:"participant_id"`
	GameID        string `json:"game_id"`
	Username      string `json:"username"`
}
