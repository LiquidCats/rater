package data

type APIRequest struct {
	Pairs []string `json:"pairs"`
}

type APIResponse struct {
	Ok   string            `json:"ok"`
	Data map[string]Ticker `json:"data"`
}

type Ticker struct {
	LastTradePrice string `json:"lastTradePrice"`
}
