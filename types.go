package idex

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

// OrderBook holds bid and ask orders for a market
type OrderBook struct {
	Bids []Order `json:"bids"`
	Asks []Order `json:"asks"`
}

// Order held in OrderBook
type Order struct {
	Price     string  `json:"price"`
	Amount    string  `json:"amount"`
	Total     string  `json:"total"`
	OrderHash string  `json:"orderHash"`
	Params    *Params `json:"params"`
}

// OpenOrder for a market or user
type OpenOrder struct {
	Timestamp   int     `json:"timestamp"`
	Price       string  `json:"price"`
	Amount      string  `json:"amount"`
	Total       string  `json:"total"`
	OrderHash   string  `json:"orderHash"`
	Market      string  `json:"market"`
	Type        string  `json:"type"`
	OrderNumber int     `json:"orderNumber"`
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
	USDValue        string `json:"usdValue"`
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

// TradeInserted event
type TradeInserted struct {
	ID              int    `json:"id"`
	Price           string `json:"price"`
	AmountPrecision string `json:"amountPrecision"`
	TotalPrecision  string `json:"totalPrecision"`
	Date            string `json:"date"`
	Timestamp       int    `json:"timestamp"`
	SellerFee       string `json:"sellerFee"`
	BuyerFee        string `json:"buyerFee"`
	Type            string `json:"type"`
	TokenBuy        string `json:"tokenBuy"`
	AmountBuy       string `json:"amountBuy"`
	TokenSell       string `json:"tokenSell"`
	AmountSell      string `json:"amountSell"`
	FeeMake         string `json:"feeMake"`
	FeeTake         string `json:"feeTake"`
	GasFee          string `json:"gasFee"`
	Buy             string `json:"buy"`
	V               int    `json:"v"`
	R               string `json:"r"`
	S               string `json:"s"`
	User            string `json:"user"`
	Sell            string `json:"sell"`
	Hash            string `json:"hash"`
	Nonce           int    `json:"nonce"`
	Amount          string `json:"amount"`
	USDValue        string `json:"usdValue"`
	GasFeeAdjusted  string `json:"gasFeeAdjusted"`
	UUID            string `json:"uuid"`
	UpdatedAt       string `json:"updatedAt"`
	CreatedAt       string `json:"createdAt"`
}

// OrderInserted event
type OrderInserted struct {
	Complete        bool   `json:"complete"`
	ID              int    `json:"id"`
	TokenBuy        string `json:"tokenBuy"`
	AmountBuy       string `json:"amountBuy"`
	TokenSell       string `json:"tokenSell"`
	AmountSell      string `json:"amountSell"`
	Expires         int    `json:"expires"`
	Nonce           int    `json:"nonce"`
	User            string `json:"user"`
	V               int    `json:"v"`
	R               string `json:"r"`
	S               string `json:"s"`
	Hash            string `json:"hash"`
	FeeDiscount     string `json:"feeDiscount"`
	RewardsMultiple string `json:"rewardsMultiple"`
	UpdatedAt       string `json:"updatedAt"`
	CreatedAt       string `json:"createdAt"`
}

// PushCancel event
type PushCancel struct {
	ID        int    `json:"id"`
	Hash      string `json:"hash"`
	User      string `json:"user"`
	V         int    `json:"v"`
	R         string `json:"r"`
	S         string `json:"s"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}

// Method name of websocket event
type Method struct {
	Method string `json:"method"`
}

// SocketResponse holds messages to pass back from the websocket
type SocketResponse struct {
	OrderInserted *OrderInserted
	TradeInserted *TradeInserted
	PushCancel    *PushCancel
	Error         error
}
