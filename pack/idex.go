package idex

import (
	"encoding/json"
	"fmt"
)

// Idex service
type Idex struct {
	Client *Client
}

const (
	url = "https://api.idex.market"
)

// New instance of an Idex
func New() *Idex {
	return &Idex{Client: NewClient(url)}
}

// Ticker for the market
func (i *Idex) Ticker(market string) (*Ticker, error) {
	if market == "" {
		return nil, fmt.Errorf("market is required")
	}

	payload := fmt.Sprintf(`{"market":"%s"}`, market)
	body, err := i.Client.do("returnTicker", payload)
	if err != nil {
		return nil, err
	}

	if string(body) == "{}" {
		return nil, fmt.Errorf("market %v not found", market)
	}

	t := &Ticker{}
	if err := json.Unmarshal(body, t); err != nil {
		return nil, err
	}

	return t, nil
}

// Tickers for all markets
func (i *Idex) Tickers() (map[string]*Ticker, error) {
	body, err := i.Client.do("returnTicker", "")
	if err != nil {
		return nil, err
	}

	ts := make(map[string]*Ticker)
	if err := json.Unmarshal(body, &ts); err != nil {
		return nil, err
	}

	return ts, nil
}

// Volume24 returns 24-hour volume for all markets
func (i *Idex) Volume24() (*Volume, error) {
	body, err := i.Client.do("return24Volume", "")
	if err != nil {
		return nil, err
	}

	vol := &Volume{}
	// using custom UnmarshalJSON method for returned structure
	if err := json.Unmarshal(body, &vol); err != nil {
		return nil, err
	}

	return vol, nil
}

// OrderBook for a market
func (i *Idex) OrderBook(market string) (*OrderBook, error) {
	if market == "" {
		return nil, fmt.Errorf("market is required")
	}

	payload := fmt.Sprintf(`{"market":"%s"}`, market)
	body, err := i.Client.do("returnOrderBook", payload)
	if err != nil {
		return nil, err
	}

	ob := &OrderBook{}
	if err := json.Unmarshal(body, &ob); err != nil {
		return nil, err
	}

	return ob, nil
}

// TODO:
// returnOpenOrders
// returnTradeHistory
// returnCurrencies
// returnBalances
// returnCompleteBalances
// returnDepositsWithdrawals
// returnOrderTrades
// returnNextNonce
// returnContractAddress

// authenticated methods:
// order
// trade
// cancel
// withdraw
