package rating

type Rating struct {
	RatingID      string  `json:"ratingId"`
	GameID        string  `json:"gameId"`
	ParticipantID string  `json:"participantId"`
	Username      string  `json:"username,omitempty"`
	WineID        string  `json:"wineId"`
	SightRating   float64 `json:"sightRating"`
	AromaRating   float64 `json:"aromaRating"`
	TasteRating   float64 `json:"tasteRating"`
	OverallRating float64 `json:"overallRating"`
	TotalRating   float64 `json:"totalRating"`
	Comments      string  `json:"comments"`
}
