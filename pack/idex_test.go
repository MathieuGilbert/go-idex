package idex

import (
	"io/ioutil"
	"net/http"
	"testing"

	mockhttp "github.com/karupanerura/go-mock-http-response"
)

// ClientMock satisfies the Poster interface, including name of fixture file.
type APIMock struct {
	Fixture string
}

func (a APIMock) Post(endpoint, payload string) ([]byte, error) {
	if a.Fixture == "" {
		return nil, nil
	}

	b, err := ioutil.ReadFile("testdata/" + a.Fixture)
	if err != nil {
		panic(err)
	}
	return b, nil
}

func mockResponse(statusCode int, headers map[string]string, body []byte) {
	http.DefaultClient = mockhttp.NewResponseMock(statusCode, headers, body).MakeClient()
}

// TestPost hits the real API.Post method,
func TestPost(t *testing.T) {
	r := "response"
	mockResponse(http.StatusOK, map[string]string{"Content-Type": "application/json"}, []byte(r))

	api := &API{URL: "https://api.example.com"}
	resp, err := api.Post("returnTicker", `{"market":"ETH_SAN"}`)
	if err != nil {
		t.Error("should not be an error")
	}
	if string(resp) != r {
		t.Error("should get the response")
	}
}

func TestNew(t *testing.T) {
	idex := New(APIMock{Fixture: ""})
	if idex == nil {
		t.Error("idex should be created")
	}
}

func TestTicker(t *testing.T) {
	idex := New(APIMock{Fixture: "ticker.json"})

	tk, err := idex.Ticker("ETH_AUC")
	if err != nil {
		t.Errorf("should not be an error: %v\n", err)
	}
	if tk == nil {
		t.Error("ticker should not be nil")
	}
	if tk.Last != "0.00555" {
		t.Errorf("last should be 0.00555, was: %v", tk.Last)
	}
}

func TestTickerBadMarket(t *testing.T) {
	idex := New(APIMock{Fixture: "empty.json"})

	_, err := idex.Ticker("INVALID")
	if err == nil {
		t.Error("should be an error")
	}

	_, err = idex.Ticker("")
	if err == nil {
		t.Error("should be an error")
	}

}

func TestTickers(t *testing.T) {
	idex := New(APIMock{Fixture: "tickers.json"})

	tks, err := idex.Tickers()
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if tks == nil {
		t.Error("tickers should not be nil")
	}
	if was, exp := len(tks), 3; was != exp {
		t.Errorf("there should be %v tickers, was: %v", exp, was)
	}
	if tks["ETH_AUC"] == nil {
		t.Error("ETH_SAN market should exist")
	}
	if was, exp := tks["ETH_AUC"].Last, "0.000207186721315823"; was != exp {
		t.Errorf("last should %v, was: %v", exp, was)
	}
}

func TestVolume24(t *testing.T) {
	idex := New(APIMock{Fixture: "volume24.json"})

	vol, err := idex.Volume24()
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if was, exp := len(vol.Markets), 3; was != exp {
		t.Errorf("there should be %v volumes, was: %v", exp, was)
	}
	if vol.Markets["ETH_AUC"] == nil {
		t.Error("market should exist")
	}
	if was, exp := vol.Markets["ETH_AUC"]["AUC"], "209849.917899637864109753"; was != exp {
		t.Errorf("volume should be %v, was: %v", exp, was)
	}
	if vol.TotalETH != "14148.11678323491238745" {
		t.Errorf("TotalETH should be 12.234, was: %v", vol.TotalETH)
	}
}

func TestReturnOrderBook(t *testing.T) {
	idex := New(APIMock{Fixture: "orderBook.json"})

	ob, err := idex.OrderBook("ETH_SAN")
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if ob.Asks == nil || len(ob.Asks) == 0 {
		t.Error("missing asks")
	}
	if was, exp := len(ob.Asks), 5; was != exp {
		t.Errorf("there should be %v asks, was: %v", exp, was)
	}
	if ob.Bids == nil || len(ob.Bids) == 0 {
		t.Error("missing bids")
	}
	if was, exp := len(ob.Bids), 6; was != exp {
		t.Errorf("there should be %v bids, was: %v", exp, was)
	}
	if was, exp := ob.Asks[0].Price, "0.00342003"; was != exp {
		t.Errorf("first ask price should be %v, was: %v", exp, was)
	}
	if was, exp := ob.Bids[0].Amount, "2222"; was != exp {
		t.Errorf("first bid amount should be %v, was: %v", exp, was)
	}
}

func TestOpenOrdersMarket(t *testing.T) {
	idex := New(APIMock{Fixture: "openOrdersMarket.json"})

	os, err := idex.OpenOrders("ETH_SAN", "")
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if was, exp := len(os), 12; was != exp {
		t.Errorf("there should be %v orders, was: %v", exp, was)
	}
	if was, exp := os[0].Timestamp, 1518721278; was != exp {
		t.Errorf("first open order timestamp should be %v, was: %v", exp, was)
	}
	if was, exp := os[0].Params.AmountBuy, "1999999999999999999"; was != exp {
		t.Errorf("first open order amount buy should be %v, was: %v", exp, was)
	}
}

func TestOpenOrdersUser(t *testing.T) {
	idex := New(APIMock{Fixture: "openOrdersUser.json"})

	os, err := idex.OpenOrders("", "0x1234567890")
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if was, exp := len(os), 10; was != exp {
		t.Errorf("there should be %v orders, was: %v", exp, was)
	}
	if was, exp := os[0].Market, "ETH_PARETO"; was != exp {
		t.Errorf("first open order market should be %v, was: %v", exp, was)
	}
	if was, exp := os[0].Params.AmountSell, "20146610957659998010000"; was != exp {
		t.Errorf("first open order market should be %v, was: %v", exp, was)
	}

	_, err = idex.OpenOrders("", "")
	if err == nil {
		t.Error("should be an error")
	}
}

func TestTradeHistoryMarket(t *testing.T) {
	idex := New(APIMock{Fixture: "tradeHistoryMarket.json"})

	ts, err := idex.TradeHistoryMarket("ETH_SAN", "", 0, 0)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if was, exp := len(ts), 22; was != exp {
		t.Errorf("there should be %v trades, was: %v", exp, was)
	}
	if was, exp := ts[0].Total, "0.150496412698412699"; was != exp {
		t.Errorf("first trade history total should be %v, was: %v", exp, was)
	}
	if was, exp := ts[1].Total, "0.5418851999999997"; was != exp {
		t.Errorf("second trade history total should be %v, was: %v", exp, was)
	}

	ts, err = idex.TradeHistoryMarket("ETH_SAN", "0x1234567890", 0, 0)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(ts) == 0 {
		t.Error("there should be at least one trade")
	}

	ts, err = idex.TradeHistoryMarket("ETH_SAN", "", 1531000000, 1532000000)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(ts) == 0 {
		t.Error("there should be at least one trade")
	}

	_, err = idex.TradeHistoryMarket("", "", 0, 0)
	if err == nil {
		t.Error("should be an error")
	}
}

func TestTradeHistoryUser(t *testing.T) {
	idex := New(APIMock{Fixture: "tradeHistoryUser.json"})

	ts, err := idex.TradeHistoryUser("0x1234567890", 0, 0)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if was, exp := len(ts), 4; was != exp {
		t.Errorf("there should be %v markets traded, was: %v", exp, was)
	}
	if was, exp := len(ts["ETH_SAN"]), 2; was != exp {
		t.Errorf("there should be %v trades on ETH_SAN, was: %v", exp, was)
	}
	if was, exp := ts["ETH_SAN"][0].Total, "0.150496412698412699"; was != exp {
		t.Errorf("first trade history total should be %v, was: %v", exp, was)
	}
	if was, exp := ts["ETH_SAN"][0].USDValue, "105.3474888888888893"; was != exp {
		t.Errorf("first trade history usdValue should be %v, was: %v", exp, was)
	}

	_, err = idex.TradeHistoryUser("", 0, 0)
	if err == nil {
		t.Error("should be an error")
	}
}

func TestCurrencies(t *testing.T) {
	idex := New(APIMock{Fixture: "currencies.json"})

	cs, err := idex.Currencies()
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if was, exp := len(cs), 3; was != exp {
		t.Errorf("there should be %v currencies, was: %v", exp, was)
	}
	if cs["1ST"] == nil {
		t.Error("1ST should be included")
	}
	if was, exp := cs["1ST"].Name, "FirstBlood"; was != exp {
		t.Errorf("1ST currency name should be %v, was %v", exp, was)
	}
}

func TestBalances(t *testing.T) {
	idex := New(APIMock{Fixture: "balances.json"})

	bs, err := idex.Balances("0x1234567890")
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(bs) == 0 {
		t.Error("there should be at least one balance")
	}
	if was, exp := len(bs), 4; was != exp {
		t.Errorf("there should be %v deposits, was: %v", exp, was)
	}
	if was, exp := bs["NPXS"], "0.000000000000055555"; was != exp {
		t.Errorf("NPXS balance should be %v, was: %v", exp, was)
	}

	_, err = idex.Balances("")
	if err == nil {
		t.Error("should be an error")
	}
}

func TestCompleteBalances(t *testing.T) {
	idex := New(APIMock{Fixture: "completeBalances.json"})

	bs, err := idex.CompleteBalances("0x1234567890")
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if was, exp := len(bs), 4; was != exp {
		t.Errorf("there should be %v balances, was: %v", exp, was)
	}
	if bs["ETH"] == nil {
		t.Error("ETH should be included")
	}
	if was, exp := bs["ETH"].Available, "0.987654321000000"; was != exp {
		t.Errorf("ETH available balance should be %v, was: %v", exp, was)
	}
	if was, exp := bs["ETH"].OnOrders, "12.777"; was != exp {
		t.Errorf("ETH on order balance should be %v, was: %v", exp, was)
	}

	_, err = idex.CompleteBalances("")
	if err == nil {
		t.Error("should be an error")
	}
}

func TestDepositsWithdrawals(t *testing.T) {
	idex := New(APIMock{Fixture: "depositsWithdrawals.json"})

	ds, ws, err := idex.DepositsWithdrawals("0x1234567890", 0, 0)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if was, exp := len(ds), 2; was != exp {
		t.Errorf("there should be %v deposits, was: %v", exp, was)
	}
	if was, exp := len(ws), 2; was != exp {
		t.Errorf("there should be %v withdrawals, was: %v", exp, was)
	}

	_, _, err = idex.DepositsWithdrawals("0x1234567890", 1510000000, 1540000000)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}

	_, _, err = idex.DepositsWithdrawals("", 0, 0)
	if err == nil {
		t.Error("should be an error")
	}
}

func TestOrderTrades(t *testing.T) {
	idex := New(APIMock{Fixture: "orderTrades.json"})

	ts, err := idex.OrderTrades("0x9876543210")
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if was, exp := len(ts), 1; was != exp {
		t.Errorf("there should be %v trades, was: %v", exp, was)
	}
	if was, exp := ts[0].Amount, "74.503174603174603704"; was != exp {
		t.Errorf("first trade's amount should be %v, was: %v", exp, was)
	}
	if was, exp := ts[0].Price, "0.00202"; was != exp {
		t.Errorf("first trade's price should be %v, was: %v", exp, was)
	}

	_, err = idex.OrderTrades("")
	if err == nil {
		t.Error("should be an error")
	}

}

func TestNextNonce(t *testing.T) {
	idex := New(APIMock{Fixture: "nextNonce.json"})

	n, err := idex.NextNonce("0x1234567890")
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if exp := 2468; n != exp {
		t.Errorf("nonce should be %v, was: %v", exp, n)
	}

	_, err = idex.NextNonce("")
	if err == nil {
		t.Error("should be an error")
	}
}

func TestContractAddress(t *testing.T) {
	idex := New(APIMock{Fixture: "contractAddress.json"})

	a, err := idex.ContractAddress()
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	// current IDEX contract address
	if a != "0x2a0c0dbecc7e4d658f48e01e3fa353f44050c208" {
		t.Errorf("address has changed: %v\n", a)
	}
}
