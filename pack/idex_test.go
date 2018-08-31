package idex

import (
	"testing"
)

const (
	skip   = true
	market = "ETH_SAN"
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
		//t.Skip()
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
