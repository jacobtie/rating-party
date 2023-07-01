package wine

type Wine struct {
	WineID   string `json:"wineId"`
	WineName string `json:"wineName"`
	WineCode string `json:"wineCode"`
	WineYear int    `json:"wineYear"`
}
