package idex

import (
	"testing"
)

const (
	skip   = true
	market = "ETH_SAN"
	user   = "0xa4e451a0e8fcb8a1c29bbe3c96b7e347a674e616"
)

func TestNew(t *testing.T) {
	idex := New()

	if idex.Client.URL != url {
		t.Error("should set the URL")
	}
}

func TestTicker(t *testing.T) {
	if skip {
		t.Skip()
	}

	idex := New()

	tk, err := idex.Ticker(market)
	if err != nil {
		t.Errorf("should not be an error: %v\n", err)
	}
	if tk == nil {
		t.Error("ticker should not be nil")
	}
}

func TestTickerBadMarket(t *testing.T) {
	if skip {
		t.Skip()
	}

	idex := New()

	tk, err := idex.Ticker("INVALID")
	if err == nil {
		t.Error("should be an error")
	}
	if tk != nil {
		t.Errorf("ticker should be nil: %v", tk)
	}

	tk, err = idex.Ticker("")
	if err == nil {
		t.Error("should be an error")
	}
	if tk != nil {
		t.Errorf("ticker should be nil: %v", tk)
	}
}

func TestTickers(t *testing.T) {
	if skip {
		t.Skip()
	}

	idex := New()

	tks, err := idex.Tickers()
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if tks == nil {
		t.Error("tickers should not be nil")
	}
	if len(tks) == 0 {
		t.Error("there should be at least one ticker")
	}
	if tks[market] == nil {
		t.Errorf("%v market should exist", market)
	}
}

func TestVolume24(t *testing.T) {
	if skip {
		t.Skip()
	}

	idex := New()

	vol, err := idex.Volume24()
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(vol.Markets) == 0 {
		t.Error("should be at least one volume")
	}
	if vol.Markets[market] == nil {
		t.Errorf("%v market should exist", market)
	}
	if vol.TotalETH == "" {
		t.Error("TotalETH should be set")
	}
}

func TestReturnOrderBook(t *testing.T) {
	if skip {
		t.Skip()
	}

	idex := New()

	ob, err := idex.OrderBook(market)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if ob.Asks == nil || len(ob.Asks) == 0 {
		t.Error("missing asks")
	}
	if ob.Bids == nil || len(ob.Bids) == 0 {
		t.Error("missing bids")
	}
}

func TestOpenOrdersMarket(t *testing.T) {
	if skip {
		t.Skip()
	}

	idex := New()

	os, err := idex.OpenOrders(market, "")
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(os) == 0 {
		t.Error("there should be at least one open order")
	}
}

func TestOpenOrdersAddress(t *testing.T) {
	if skip {
		t.Skip()
	}

	idex := New()

	os, err := idex.OpenOrders("", user)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(os) == 0 {
		t.Error("there should be at least one open order")
	}

	_, err = idex.OpenOrders("", "")
	if err == nil {
		t.Error("should be an error")
	}
}

func TestTradeHistoryMarket(t *testing.T) {
	if skip {
		t.Skip()
	}

	idex := New()

	ts, err := idex.TradeHistoryMarket(market, "", 0, 0)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(ts) == 0 {
		t.Error("there should be at least one trade")
	}

	ts, err = idex.TradeHistoryMarket(market, user, 0, 0)
	if err != nil {
		// dependent on user's trade history
		t.Errorf("should not be an error: %v", err)
	}
	if len(ts) == 0 {
		t.Error("there should be at least one trade")
	}

	ts, err = idex.TradeHistoryMarket(market, "", 1531000000, 1532000000)
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
	if skip {
		t.Skip()
	}

	idex := New()

	ts, err := idex.TradeHistoryUser(user, 0, 0)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(ts) == 0 {
		// dependent on user's trade history
		t.Error("there should be at least one trade")
	}

	_, err = idex.TradeHistoryUser("", 0, 0)
	if err == nil {
		t.Error("should be an error")
	}
}

func TestCurrencies(t *testing.T) {
	if skip {
		t.Skip()
	}

	idex := New()

	cs, err := idex.Currencies()
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(cs) == 0 {
		t.Error("there should be at least one currency")
	}
	if cs["ETH"] == nil {
		t.Error("ETH should be included")
	}
}

func TestBalances(t *testing.T) {
	if skip {
		t.Skip()
	}

	idex := New()

	bs, err := idex.Balances(user)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(bs) == 0 {
		t.Error("there should be at least one balance")
	}
	if bs["DNA"] == "" {
		// DNA chosen by looking up user's balance
		t.Errorf("DNA should be included: %v", bs)
	}
}

func TestCompleteBalances(t *testing.T) {
	if skip {
		t.Skip()
	}

	idex := New()

	bs, err := idex.CompleteBalances(user)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(bs) == 0 {
		t.Error("there should be at least one balance")
	}
	if bs["DNA"] == nil {
		// DNA chosen by looking up user's balance
		t.Errorf("DNA should be included: %v", bs)
	}
}

func TestDepositsWithdrawals(t *testing.T) {
	if skip {
		t.Skip()
	}

	idex := New()

	ds, ws, err := idex.DepositsWithdrawals(user, 0, 0)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(ds) == 0 {
		t.Error("there should be at least one deposit")
	}
	if len(ws) == 0 {
		t.Error("there should be at least one withdrawal")
	}

	ds, ws, err = idex.DepositsWithdrawals(user, 1510000000, 1540000000)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(ds) == 0 {
		t.Error("there should be at least one deposit")
	}
	if len(ws) == 0 {
		t.Error("there should be at least one withdrawal")
	}

	_, _, err = idex.DepositsWithdrawals("", 0, 0)
	if err == nil {
		t.Error("should be an error")
	}
}

func TestOrderTrades(t *testing.T) {
	if skip {
		t.Skip()
	}

	idex := New()

	// found by searching through ETH_SAN trade histories
	hash := "0xfb048cab1b72e002b351962927d8f6b3d538f344b71615220d9f8eebde8f735f"

	ts, err := idex.OrderTrades(hash)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(ts) == 0 {
		t.Error("there should be at least one associated trade")
	}

	ts, err = idex.OrderTrades("0x1234")
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if len(ts) != 0 {
		t.Error("there should not be any trades")
	}

	_, err = idex.OrderTrades("")
	if err == nil {
		t.Error("should be an error")
	}

}

func TestNextNonce(t *testing.T) {
	if skip {
		t.Skip()
	}

	idex := New()

	n, err := idex.NextNonce(user)
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	if n == 0 {
		t.Error("nonce should be set")
	}

	_, err = idex.NextNonce("")
	if err == nil {
		t.Error("should be an error")
	}
}

func TestContractAddress(t *testing.T) {
	if skip {
		t.Skip()
	}

	idex := New()

	a, err := idex.ContractAddress()
	if err != nil {
		t.Errorf("should not be an error: %v", err)
	}
	// current IDEX contract address
	if a != "0x2a0c0dbecc7e4d658f48e01e3fa353f44050c208" {
		t.Errorf("address has changed: %v\n", a)
	}
}
