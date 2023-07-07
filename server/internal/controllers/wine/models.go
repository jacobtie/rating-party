package wine

type Wine struct {
	WineID   string `json:"wineId"`
	WineName string `json:"wineName,omitempty"`
	WineCode string `json:"wineCode"`
	WineYear int    `json:"wineYear,omitempty"`
}
