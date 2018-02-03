package v2

import (
	"encoding/json"
	"fmt"
)

// TradeStructure holds executed trade information
type Trade struct {
	ID     int64
	MTS    int64
	Amount float64
	Price  float64
}

func (ts Trade) UnmarshalJSON(b []byte) error {
	tmp := []interface{}{&ts.ID, &ts.MTS, &ts.Amount, &ts.Price}
	wantLen := len(tmp)
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Trade: %d != %d", g, e)
	}
	return nil
}
