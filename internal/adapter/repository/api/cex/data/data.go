package data

type ApiRequest struct {
	Pairs []string `json:"pairs"`
}

type ApiResponse struct {
	Ok   string            `json:"ok"`
	Data map[string]Ticker `json:"data"`
}

type Ticker struct {
	LastTradePrice string `json:"lastTradePrice"`
}
