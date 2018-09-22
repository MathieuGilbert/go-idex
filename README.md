### TODO

  - Contract-backed trade functions:
    - order
    - trade
    - cancel
    - withdraw


### REST Example

```
func Ticker(mkt string) {
    i := idex.New()
	t, _ := i.API.Ticker("ETH_AUC")
    fmt.Printf("ticker: %+v\n", t)
}

```

### Websocket Example

```
func ConsumeIdex(ctx context.Context) {
	i := idex.New()
	if err := i.Socket.Connect(); err != nil {
		log.Panic(err)
	}
	defer i.Socket.Conn.Close()

	response := make(chan idex.SocketResponse)
	go i.Socket.Monitor(response)

	for {
		select {
		case r := <-response:
			if r.PushCancel != nil {
				fmt.Printf("got a cancel: %+v\n", r.PushCancel)
			}
			if r.TradeInserted != nil {
				fmt.Printf("got a trade: %+v\n", r.TradeInserted)
			}
			if r.OrderInserted != nil {
				fmt.Printf("got an order: %+v\n", r.OrderInserted)
			}
			if r.Error != nil {
				fmt.Printf("got an error: %+v\n", r.Error)
			}
		case <-ctx.Done():
			err := i.Socket.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error writing close message:", err)
				return
			}

			<-time.After(time.Second)
			return
		}
	}
}
```
