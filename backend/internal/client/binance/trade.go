package binance

import (
	"encoding/json"
	"time"
)

// Trade is a single execution from the @trade stream.
type Trade struct {
	EventType    string `json:"e"`
	EventTime    int64  `json:"E"`
	Symbol       string `json:"s"`
	TradeID      int64  `json:"t"`
	Price        string `json:"p"`
	Quantity     string `json:"q"`
	TradeTime    int64  `json:"T"`
	IsBuyerMaker bool   `json:"m"`
}

// Time is when the trade executed, not when Binance published it.
func (t Trade) Time() time.Time {
	return time.UnixMilli(t.TradeTime)
}

// streamMessage is the combined-stream envelope.
type streamMessage struct {
	Stream string          `json:"stream"`
	Data   json.RawMessage `json:"data"`
}
