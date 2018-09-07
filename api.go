package idex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// API for requests
type API struct {
	URL string
}

// Post returns the result of a POST to the endpoint with the payload
func (a *API) Post(endpoint, payload string) ([]byte, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", a.URL, endpoint), bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// Ticker for the market
func (a *API) Ticker(market string) (t *Ticker, err error) {
	if market == "" {
		err = fmt.Errorf("market is required")
		return
	}

	payload := fmt.Sprintf(`{"market":"%s"}`, market)

	body, err := a.Post("returnTicker", payload)
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
func (a *API) Tickers() (t map[string]*Ticker, err error) {
	body, err := a.Post("returnTicker", "")
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &t)
	return
}

// Volume24 returns 24-hour volume for all markets
func (a *API) Volume24() (v *Volume, err error) {
	body, err := a.Post("return24Volume", "")
	if err != nil {
		return
	}

	// using custom UnmarshalJSON method for returned structure
	err = json.Unmarshal(body, &v)
	return
}

// OrderBook for a market
func (a *API) OrderBook(market string) (ob *OrderBook, err error) {
	if market == "" {
		err = fmt.Errorf("market is required")
		return
	}

	payload := fmt.Sprintf(`{"market":"%s"}`, market)

	body, err := a.Post("returnOrderBook", payload)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ob)
	return
}

// OpenOrders all open orders for market and/or user address
func (a *API) OpenOrders(market, address string) (os []*OpenOrder, err error) {
	if market == "" && address == "" {
		err = fmt.Errorf("market or address is required")
		return
	}

	payload := fmt.Sprintf(`{"market":"%s", "address":"%s"}`, market, address)

	body, err := a.Post("returnOpenOrders", payload)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &os)
	return
}

// TradeHistoryMarket trade history for a market, filterable by user and timestamps
// API limited to 200 trades
func (a *API) TradeHistoryMarket(market, address string, start, end int) (ts []*Trade, err error) {
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

	body, err := a.Post("returnTradeHistory", payload)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ts)
	return
}

// TradeHistoryUser trade history for a user across all markets, filterable by timestamps
// API limited to 200 trades
func (a *API) TradeHistoryUser(address string, start, end int) (ts map[string][]*Trade, err error) {
	if address == "" {
		err = fmt.Errorf("address is required")
		return
	}

	payload := fmt.Sprintf(`{"address":"%s", "start":%d, "end":%d}`, address, start, end)

	body, err := a.Post("returnTradeHistory", payload)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ts)
	return
}

// Currencies returns all supported currencies
func (a *API) Currencies() (cs map[string]*Currency, err error) {
	body, err := a.Post("returnCurrencies", "")
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &cs)
	return
}

// Balances returns available balances
func (a *API) Balances(address string) (bs map[string]string, err error) {
	if address == "" {
		err = fmt.Errorf("address is required")
		return
	}

	payload := fmt.Sprintf(`{"address":"%s"}`, address)

	body, err := a.Post("returnBalances", payload)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &bs)
	return
}

// CompleteBalances returns available balances with balance in orders
func (a *API) CompleteBalances(address string) (bs map[string]*Balance, err error) {
	if address == "" {
		err = fmt.Errorf("address is required")
		return
	}

	payload := fmt.Sprintf(`{"address":"%s"}`, address)

	body, err := a.Post("returnCompleteBalances", payload)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &bs)
	return
}

// DepositsWithdrawals returns the user's deposits and withdrawals
func (a *API) DepositsWithdrawals(address string, start, end int) (ds []*Deposit, ws []*Withdrawal, err error) {
	if address == "" {
		err = fmt.Errorf("address is required")
		return
	}

	payload := fmt.Sprintf(`{"address":"%s", "start":%d, "end":%d}`, address, start, end)

	body, err := a.Post("returnDepositsWithdrawals", payload)
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
func (a *API) OrderTrades(hash string) (ts []*Trade, err error) {
	if hash == "" {
		err = fmt.Errorf("hash is required")
		return
	}

	payload := fmt.Sprintf(`{"orderHash":"%s"}`, hash)

	body, err := a.Post("returnOrderTrades", payload)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ts)
	return
}

// NextNonce returns the next available nonce
func (a *API) NextNonce(address string) (nonce int, err error) {
	if address == "" {
		err = fmt.Errorf("address is required")
		return
	}

	payload := fmt.Sprintf(`{"address":"%s"}`, address)

	body, err := a.Post("returnNextNonce", payload)
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
func (a *API) ContractAddress() (address string, err error) {
	body, err := a.Post("returnContractAddress", "")
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
