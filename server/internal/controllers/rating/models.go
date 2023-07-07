package rating

type Rating struct {
	RatingID      string `json:"ratingId"`
	GameID        string `json:"gameId"`
	ParticipantID string `json:"participantId"`
	WineID        string `json:"wineId"`
	SightRating   int    `json:"sightRating"`
	AromaRating   int    `json:"aromaRating"`
	TasteRating   int    `json:"tasteRating"`
	OverallRating int    `json:"overallRating"`
	TotalRating   int    `json:"totalRating"`
	Comments      string `json:"comments"`
}
