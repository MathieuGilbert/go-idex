package idex

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Idex service
type Idex struct {
	API Poster
}

// Poster does a call to an endpoint.
// It's here so tests can skip the HTTP call.
type Poster interface {
	Post(endpoint, payload string) ([]byte, error)
}

const (
	url = "https://api.idex.market"
)

// New instance of an Idex
func New(a Poster) *Idex {
	return &Idex{API: a}
}

// Ticker for the market
func (i *Idex) Ticker(market string) (t *Ticker, err error) {
	if market == "" {
		err = fmt.Errorf("market is required")
		return
	}

	payload := fmt.Sprintf(`{"market":"%s"}`, market)

	body, err := i.API.Post("returnTicker", payload)
	if err != nil {
		return
	}

	if strings.TrimSpace(string(body)) == "{}" {
		err = fmt.Errorf("market %v not found", market)
		return
	}

	err = json.Unmarshal(body, &t)
	return
}

// Tickers for all markets
func (i *Idex) Tickers() (t map[string]*Ticker, err error) {
	body, err := i.API.Post("returnTicker", "")
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &t)
	return
}

// Volume24 returns 24-hour volume for all markets
func (i *Idex) Volume24() (v *Volume, err error) {
	body, err := i.API.Post("return24Volume", "")
	if err != nil {
		return
	}

	// using custom UnmarshalJSON method for returned structure
	err = json.Unmarshal(body, &v)
	return
}

// OrderBook for a market
func (i *Idex) OrderBook(market string) (ob *OrderBook, err error) {
	if market == "" {
		err = fmt.Errorf("market is required")
		return
	}

	payload := fmt.Sprintf(`{"market":"%s"}`, market)

	body, err := i.API.Post("returnOrderBook", payload)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ob)
	return
}

// OpenOrders all open orders for market and/or user address
func (i *Idex) OpenOrders(market, address string) (os []*OpenOrder, err error) {
	if market == "" && address == "" {
		err = fmt.Errorf("market or address is required")
		return
	}

	payload := fmt.Sprintf(`{"market":"%s", "address":"%s"}`, market, address)

	body, err := i.API.Post("returnOpenOrders", payload)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &os)
	return
}

// TradeHistoryMarket trade history for a market, filterable by user and timestamps
// API limited to 200 trades
func (i *Idex) TradeHistoryMarket(market, address string, start, end int) (ts []*Trade, err error) {
	if market == "" {
		err = fmt.Errorf("market is required")
		return
	}

	var payload string
	if address == "" {
		payload = fmt.Sprintf(`{"market":"%s", "start":%d, "end":%d}`, market, start, end)
	} else {
		payload = fmt.Sprintf(`{"market":"%s", "address":"%s", "start":%d, "end":%d}`, market, address, start, end)
	}

	body, err := i.API.Post("returnTradeHistory", payload)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ts)
	return
}

// TradeHistoryUser trade history for a user across all markets, filterable by timestamps
// API limited to 200 trades
func (i *Idex) TradeHistoryUser(address string, start, end int) (ts map[string][]*Trade, err error) {
	if address == "" {
		err = fmt.Errorf("address is required")
		return
	}

	payload := fmt.Sprintf(`{"address":"%s", "start":%d, "end":%d}`, address, start, end)

	body, err := i.API.Post("returnTradeHistory", payload)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ts)
	return
}

// Currencies returns all supported currencies
func (i *Idex) Currencies() (cs map[string]*Currency, err error) {
	body, err := i.API.Post("returnCurrencies", "")
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &cs)
	return
}

// Balances returns available balances
func (i *Idex) Balances(address string) (bs map[string]string, err error) {
	if address == "" {
		err = fmt.Errorf("address is required")
		return
	}

	payload := fmt.Sprintf(`{"address":"%s"}`, address)

	body, err := i.API.Post("returnBalances", payload)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &bs)
	return
}

// CompleteBalances returns available balances with balance in orders
func (i *Idex) CompleteBalances(address string) (bs map[string]*Balance, err error) {
	if address == "" {
		err = fmt.Errorf("address is required")
		return
	}

	payload := fmt.Sprintf(`{"address":"%s"}`, address)

	body, err := i.API.Post("returnCompleteBalances", payload)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &bs)
	return
}

// DepositsWithdrawals returns the user's deposits and withdrawals
func (i *Idex) DepositsWithdrawals(address string, start, end int) (ds []*Deposit, ws []*Withdrawal, err error) {
	if address == "" {
		err = fmt.Errorf("address is required")
		return
	}

	payload := fmt.Sprintf(`{"address":"%s", "start":%d, "end":%d}`, address, start, end)

	body, err := i.API.Post("returnDepositsWithdrawals", payload)
	if err != nil {
		return
	}

	type response struct {
		Deposits    []*Deposit    `json:"deposits"`
		Withdrawals []*Withdrawal `json:"withdrawals"`
	}
	r := response{}

	if err = json.Unmarshal(body, &r); err != nil {
		return
	}
	ds = r.Deposits
	ws = r.Withdrawals

	return
}

// OrderTrades returns all trades involved in the order hash
func (i *Idex) OrderTrades(hash string) (ts []*Trade, err error) {
	if hash == "" {
		err = fmt.Errorf("hash is required")
		return
	}

	payload := fmt.Sprintf(`{"orderHash":"%s"}`, hash)

	body, err := i.API.Post("returnOrderTrades", payload)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ts)
	return
}

// NextNonce returns the next available nonce
func (i *Idex) NextNonce(address string) (nonce int, err error) {
	if address == "" {
		err = fmt.Errorf("address is required")
		return
	}

	payload := fmt.Sprintf(`{"address":"%s"}`, address)

	body, err := i.API.Post("returnNextNonce", payload)
	if err != nil {
		return
	}

	type response struct {
		Nonce int `json:"nonce"`
	}
	r := response{}

	if err = json.Unmarshal(body, &r); err != nil {
		return
	}
	nonce = r.Nonce

	return
}

// ContractAddress returns the IDEX contract address
func (i *Idex) ContractAddress() (address string, err error) {
	body, err := i.API.Post("returnContractAddress", "")
	if err != nil {
		return
	}

	type response struct {
		Address string `json:"address"`
	}
	r := response{}

	if err = json.Unmarshal(body, &r); err != nil {
		return
	}
	address = r.Address

	return
}
