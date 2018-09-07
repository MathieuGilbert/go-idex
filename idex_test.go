package idex

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	mockhttp "github.com/karupanerura/go-mock-http-response"
)

func fileBytes(fileName string) []byte {
	if fileName == "" {
		panic(fmt.Errorf("file name required"))
	}

	b, err := ioutil.ReadFile("testdata/" + fileName)
	if err != nil {
		panic(err)
	}
	return b
}

func mockResponse(statusCode int, body []byte) {
	h := map[string]string{"Content-Type": "application/json"}
	http.DefaultClient = mockhttp.NewResponseMock(statusCode, h, body).MakeClient()
}

func TestNew(t *testing.T) {
	idex := New()
	if idex.API == nil {
		t.Error("idex should have an API")
	}
	if idex.Socket == nil {
		t.Error("idex should have a Socket")
	}
}

func TestTicker(t *testing.T) {
	mockResponse(http.StatusOK, fileBytes("ticker.json"))
	idex := New()

	tk, err := idex.API.Ticker("ETH_AUC")
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
	mockResponse(http.StatusOK, fileBytes("empty.json"))
	idex := New()

	_, err := idex.API.Ticker("INVALID")
	if err == nil {
		t.Error("should be an error")
	}

	_, err = idex.API.Ticker("")
	if err == nil {
		t.Error("should be an error")
	}

}

func TestTickers(t *testing.T) {
	mockResponse(http.StatusOK, fileBytes("tickers.json"))
	idex := New()

	tks, err := idex.API.Tickers()
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
	mockResponse(http.StatusOK, fileBytes("volume24.json"))
	idex := New()

	vol, err := idex.API.Volume24()
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
	mockResponse(http.StatusOK, fileBytes("orderBook.json"))
	idex := New()

	ob, err := idex.API.OrderBook("ETH_SAN")
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
	mockResponse(http.StatusOK, fileBytes("openOrdersMarket.json"))
	idex := New()

	os, err := idex.API.OpenOrders("ETH_SAN", "")
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
	mockResponse(http.StatusOK, fileBytes("openOrdersUser.json"))
	idex := New()

	os, err := idex.API.OpenOrders("", "0x1234567890")
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

	_, err = idex.API.OpenOrders("", "")
	if err == nil {
		t.Error("should be an error")
	}
}

func TestTradeHistoryMarket(t *testing.T) {
	mockResponse(http.StatusOK, fileBytes("tradeHistoryMarket.json"))
	idex := New()

	ts, err := idex.API.TradeHistoryMarket("ETH_SAN", "", 0, 0)
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

	ts, err = idex.API.TradeHistoryMarket("ETH_SAN", "0x1234567890", 0, 0)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(ts) == 0 {
		t.Error("there should be at least one trade")
	}

	ts, err = idex.API.TradeHistoryMarket("ETH_SAN", "", 1531000000, 1532000000)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(ts) == 0 {
		t.Error("there should be at least one trade")
	}

	_, err = idex.API.TradeHistoryMarket("", "", 0, 0)
	if err == nil {
		t.Error("should be an error")
	}
}

func TestTradeHistoryUser(t *testing.T) {
	mockResponse(http.StatusOK, fileBytes("tradeHistoryUser.json"))
	idex := New()

	ts, err := idex.API.TradeHistoryUser("0x1234567890", 0, 0)
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

	_, err = idex.API.TradeHistoryUser("", 0, 0)
	if err == nil {
		t.Error("should be an error")
	}
}

func TestCurrencies(t *testing.T) {
	mockResponse(http.StatusOK, fileBytes("currencies.json"))
	idex := New()

	cs, err := idex.API.Currencies()
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
	mockResponse(http.StatusOK, fileBytes("balances.json"))
	idex := New()

	bs, err := idex.API.Balances("0x1234567890")
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

	_, err = idex.API.Balances("")
	if err == nil {
		t.Error("should be an error")
	}
}

func TestCompleteBalances(t *testing.T) {
	mockResponse(http.StatusOK, fileBytes("completeBalances.json"))
	idex := New()

	bs, err := idex.API.CompleteBalances("0x1234567890")
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

	_, err = idex.API.CompleteBalances("")
	if err == nil {
		t.Error("should be an error")
	}
}

func TestDepositsWithdrawals(t *testing.T) {
	mockResponse(http.StatusOK, fileBytes("depositsWithdrawals.json"))
	idex := New()

	ds, ws, err := idex.API.DepositsWithdrawals("0x1234567890", 0, 0)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if was, exp := len(ds), 2; was != exp {
		t.Errorf("there should be %v deposits, was: %v", exp, was)
	}
	if was, exp := len(ws), 2; was != exp {
		t.Errorf("there should be %v withdrawals, was: %v", exp, was)
	}

	_, _, err = idex.API.DepositsWithdrawals("0x1234567890", 1510000000, 1540000000)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}

	_, _, err = idex.API.DepositsWithdrawals("", 0, 0)
	if err == nil {
		t.Error("should be an error")
	}
}

func TestOrderTrades(t *testing.T) {
	mockResponse(http.StatusOK, fileBytes("orderTrades.json"))
	idex := New()

	ts, err := idex.API.OrderTrades("0x9876543210")
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

	_, err = idex.API.OrderTrades("")
	if err == nil {
		t.Error("should be an error")
	}

}

func TestNextNonce(t *testing.T) {
	mockResponse(http.StatusOK, fileBytes("nextNonce.json"))
	idex := New()

	n, err := idex.API.NextNonce("0x1234567890")
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if exp := 2468; n != exp {
		t.Errorf("nonce should be %v, was: %v", exp, n)
	}

	_, err = idex.API.NextNonce("")
	if err == nil {
		t.Error("should be an error")
	}
}

func TestContractAddress(t *testing.T) {
	mockResponse(http.StatusOK, fileBytes("contractAddress.json"))
	idex := New()

	a, err := idex.API.ContractAddress()
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	// current IDEX contract address
	if a != "0x2a0c0dbecc7e4d658f48e01e3fa353f44050c208" {
		t.Errorf("address has changed: %v\n", a)
	}
}

func TestProcessTradesInserted(t *testing.T) {
	c := make(chan SocketResponse)
	go processTradesInserted(fileBytes("notifyTradesInserted.json"), c)

	r := <-c
	if r.Error != nil {
		t.Errorf("should not be an error: %v", r.Error)
	}
	if r.TradeInserted == nil {
		t.Error("should be a TradeInserted")
	}
	if was, exp := r.TradeInserted.Price, "0.000002117466563483"; was != exp {
		t.Errorf("inserted trade price should be %v, was: %v", exp, was)
	}
	if was, exp := r.TradeInserted.V, 27; was != exp {
		t.Errorf("inserted trade V should be %v, was: %v", exp, was)
	}

	r = <-c
	if r.Error != nil {
		t.Errorf("should not be an error: %v", r.Error)
	}
	if r.TradeInserted == nil {
		t.Error("should be a TradeInserted")
	}
	if was, exp := r.TradeInserted.Price, "0.000002118000000000"; was != exp {
		t.Errorf("inserted trade price should be %v, was: %v", exp, was)
	}
	// this case tests for when V comes back as a string instead of an int
	if was, exp := r.TradeInserted.V, 28; was != exp {
		t.Errorf("inserted trade V should be %v, was: %v", exp, was)
	}
}

func TestProcessOrderInserted(t *testing.T) {
	c := make(chan SocketResponse)
	go processOrderInserted(fileBytes("notifyOrderInserted.json"), c)

	r := <-c
	if r.Error != nil {
		t.Errorf("should not be an error: %v", r.Error)
	}
	if r.OrderInserted == nil {
		t.Error("should be an OrderInserted")
	}
	if was, exp := r.OrderInserted.TokenSell, "0x3f06b5d78406cd97bdf10f5c420b241d32759c80"; was != exp {
		t.Errorf("inserted order token sell should be %v, was: %v", exp, was)
	}
	if was, exp := r.OrderInserted.V, 27; was != exp {
		t.Errorf("inserted order V should be %v, was: %v", exp, was)
	}

	// this file has V as a string
	go processOrderInserted(fileBytes("notifyOrderInserted2.json"), c)

	r = <-c
	if r.Error != nil {
		t.Errorf("should not be an error: %v", r.Error)
	}
	if was, exp := r.OrderInserted.V, 26; was != exp {
		t.Errorf("inserted order V should be %v, was: %v", exp, was)
	}
}

func TestProcessPushCancel(t *testing.T) {
	c := make(chan SocketResponse)
	go processPushCancel(fileBytes("notifyPushCancel.json"), c)

	r := <-c
	if r.Error != nil {
		t.Errorf("should not be an error: %v", r.Error)
	}
	if r.PushCancel == nil {
		t.Error("should be a PushCancel")
	}
	if was, exp := r.PushCancel.Hash, "0x216a8e0de8c3fc08279e0fccee5b9da7011312dab4b740288729f4f77497cbaa"; was != exp {
		t.Errorf("cancel hash should be %v, was: %v", exp, was)
	}
	if was, exp := r.PushCancel.V, 28; was != exp {
		t.Errorf("cancel V should be %v, was: %v", exp, was)
	}

	// this file has V as a string
	go processPushCancel(fileBytes("notifyPushCancel2.json"), c)

	r = <-c
	if r.Error != nil {
		t.Errorf("should not be an error: %v", r.Error)
	}
	if r.PushCancel == nil {
		t.Error("should be a PushCancel")
	}
	if was, exp := r.PushCancel.V, 29; was != exp {
		t.Errorf("cancel V should be %v, was: %v", exp, was)
	}
}

func TestProcessPushCancels(t *testing.T) {
	c := make(chan SocketResponse)
	go processPushCancels(fileBytes("notifyPushCancels.json"), c)

	r := <-c
	if r.Error != nil {
		t.Errorf("should not be an error: %v", r.Error)
	}
	if r.PushCancel == nil {
		t.Error("should be a PushCancel")
	}
	if was, exp := r.PushCancel.Hash, "0xef464f5d2bd68459be5c4f16d6d34e79c9079aa61fc8b27bdfc3efa6541c2a2d"; was != exp {
		t.Errorf("cancel hash should be %v, was: %v", exp, was)
	}

	r = <-c
	if r.Error != nil {
		t.Errorf("should not be an error: %v", r.Error)
	}
	if r.PushCancel == nil {
		t.Error("should be a PushCancel")
	}
	if was, exp := r.PushCancel.Hash, "0xff464f5d2bd68459be5c4f16d6d34e79c9079aa61fc8b27bdfc3efa6541c2a2d"; was != exp {
		t.Errorf("cancel hash should be %v, was: %v", exp, was)
	}

}
