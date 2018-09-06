### TODO

#### Contract-backed trade functions
  - order
  - trade
  - cancel
  - withdraw

#### Websockets

  - handshake
  - notifyTradesInserted
  - notifyOrderInserted
  - pushCancel?
  - pushCancels?


### Example

```
func ConsumeIdex() {
	// cleanly shutdown on interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	url := "wss://v1.idex.market"
	s := idex.NewSocket(url)
	if err := s.Connect(); err != nil {
		log.Panic(err)
	}
	defer s.Conn.Close()

	if err := s.Handshake(); err != nil {
		log.Panic(err)
	}

	done := make(chan struct{})
	response := make(chan idex.SocketResponse)
	go s.Monitor(done, response)

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
				close(done)
			}
		case <-done:
			return
		case <-interrupt:
			fmt.Println("interrupted")

			err := s.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Panicf("Error writing close message: %v\n", err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
```
