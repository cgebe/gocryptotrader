package v2

// TradeStructure holds executed trade information
type Trade struct {
	ID     int64   `json:"id"`
	MTS    int64   `json:"mts"`
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}
