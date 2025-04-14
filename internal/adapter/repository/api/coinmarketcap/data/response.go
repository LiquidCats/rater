package data

type APIResponse struct {
	Data Conversion `json:"data"`
}

type Conversion struct {
	Name   string           `json:"name"`
	Symbol string           `json:"symbol"`
	Quote  map[string]Price `json:"quote"`
}

type Price struct {
	Price float64 `json:"price"`
}
