package idex

import (
	"encoding/json"
)

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

// OpenOrder in orderbook
type OpenOrder struct {
	Timestamp   int     `json:"timestamp"`
	OrderHash   string  `json:"orderHash"`
	Market      string  `json:"market"`
	Type        string  `json:"type"`
	OrderNumber int     `json:"orderNumber"`
	Price       string  `json:"price"`
	Amount      string  `json:"amount"`
	Total       string  `json:"total"`
	Params      *Params `json:"params"`
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

// Trade holds details about a trade
type Trade struct {
	Date            string `json:"date"`
	Amount          string `json:"amount"`
	Type            string `json:"type"`
	Total           string `json:"total"`
	Price           string `json:"price"`
	OrderHash       string `json:"orderHash"`
	UUID            string `json:"uuid"`
	BuyerFee        string `json:"buyerFee"`
	SellerFee       string `json:"sellerFee"`
	GasFee          string `json:"gasFee"`
	Timestamp       int    `json:"timestamp"`
	Maker           string `json:"maker"`
	Taker           string `json:"taker"`
	TransactionHash string `json:"transactionHash"`
	UsdValue        string `json:"usdValue"`
}

// Currency holds details about supported currencies
type Currency struct {
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
	Address  string `json:"address"`
}

// Balance of token available and in open orders
type Balance struct {
	Available string `json:"available"`
	OnOrders  string `json:"onOrders"`
}

// Deposit holds information about a user's deposits
type Deposit struct {
	DepositNumber   int    `json:"depositNumber"`
	Currency        string `json:"currency"`
	Amount          string `json:"amount"`
	Timestamp       int    `json:"timestamp"`
	TransactionHash string `json:"transactionHash"`
}

// Withdrawal holds information about a user's withdrawals
type Withdrawal struct {
	WithdrawalNumber int    `json:"depositNumber"`
	Currency         string `json:"currency"`
	Amount           string `json:"amount"`
	Timestamp        int    `json:"timestamp"`
	TransactionHash  string `json:"transactionHash"`
	Status           string `json:"status"`
}

// Volume maps markets to ETH and TOKEN amounts, with total eth volume
type Volume struct {
	Markets  map[string]map[string]string
	TotalETH string
}

// UnmarshalJSON custom for Volume
func (v *Volume) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &v.Markets); err != nil {
		// expecting error on totalETH field having a string value
		if !UnmarshalErrorOnType(err, "string") {
			return err
		}
	}
	delete(v.Markets, "totalETH")

	type total struct {
		TotalETH string `json:"totalETH"`
	}
	t := total{}

	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}
	v.TotalETH = t.TotalETH

	return nil
}
