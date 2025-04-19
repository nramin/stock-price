package structs

type StockPrices struct {
	Stocks     []Stock `json:"stocks,omitempty"`
	TotalValue float64 `json:"totalValue,omitempty"`
	Success    *bool   `json:"success,omitempty"`
	Error      string  `json:"error,omitempty"`
}

type Stock struct {
	Symbol   string  `json:"symbol,omitempty"`
	Price    float64 `json:"price,omitempty"`
	Quantity float64 `json:"quantity,omitempty"`
	Value    float64 `json:"value,omitempty"`
}
