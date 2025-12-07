package steam

type GetPlayerCountResponse struct {
	Data GetPlayerCountData `json:"response"`
}

type GetPlayerCountData struct {
	PlayerCount int `json:"player_count"`
}

type GetAppDetailsResponse map[string]struct {
	Data GetAppDetailsData `json:"data"`
}

type GetAppDetailsData struct {
	Name          string                     `json:"name"`
	Packages      []int                      `json:"packages"`
	PriceOverview GetAppDetailsPriceOverview `json:"price_overview"`
}

type GetAppDetailsPriceOverview struct {
	Initial  int    `json:"initial"`
	Final    int    `json:"final"`
	Currency string `json:"currency"`
}
