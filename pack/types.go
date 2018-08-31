package idex

import "encoding/json"

// Ticker data
type Ticker struct {
	Last          string `json:"last"`
	High          string `json:"high"`
	Low           string `json:"low"`
	LowestAsk     string `json:"lowestAsk"`
	HighestBid    string `json:"highestBid"`
	PercentChange string `json:"percentChange"`
	BaseVolume    string `json:"baseVolume"`
	QuoteVolume   string `json:"quoteVolume"`
}

// Order in orderbook
type Order struct {
	Price     string  `json:"price"`
	Amount    string  `json:"amount"`
	Total     string  `json:"total"`
	OrderHash string  `json:"orderHash"`
	Params    *Params `json:"params"`
}

// Params of an Order
type Params struct {
	TokenBuy      string `json:"tokenBuy"`
	BuySymbol     string `json:"buySymbol"`
	BuyPrecision  int    `json:"buyPrecision"`
	AmountBuy     string `json:"amountBuy"`
	TokenSell     string `json:"tokenSell"`
	SellSymbol    string `json:"sellSymbol"`
	SellPrecision int    `json:"sellPrecision"`
	AmountSell    string `json:"amountSell"`
	Expires       int    `json:"expires"`
	Nonce         int    `json:"nonce"`
	User          string `json:"user"`
}

// OrderBook holds bid and ask orders
type OrderBook struct {
	Bids []Order `json:"bids"`
	Asks []Order `json:"asks"`
}

// Volume maps markets to ETH and TOKEN amounts, with total eth volume
type Volume struct {
	Markets  map[string]map[string]string
	TotalETH string `json:"totalETH"`
}

// UnmarshalJSON custom for Volume
func (v *Volume) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &v.Markets); err != nil {
		switch e := err.(type) {
		case *json.UnmarshalTypeError:
			if e.Value != "string" {
				return err
			}
		default:
			return err
		}
	}
	delete(v.Markets, "totalETH")

	type total struct {
		TotalETH string `json:"totalETH"`
	}
	t := total{}

	if err := json.Unmarshal(b, &t); err != nil {
		switch e := err.(type) {
		case *json.UnmarshalTypeError:
			if e.Value != "map[string]map[string]string" {
				return err
			}
		default:
			return err
		}
	}
	v.TotalETH = t.TotalETH

	return nil
}
