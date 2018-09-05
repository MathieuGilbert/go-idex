package idex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Socket wraps the websocket connection
type Socket struct {
	Conn *websocket.Conn
	URL  string
}

// NewSocket returns a new Socket
func NewSocket(url string) *Socket {
	return &Socket{URL: url}
}

// Connect to websocket
func (s *Socket) Connect() error {
	c, resp, err := websocket.DefaultDialer.Dial(s.URL, nil)
	if err == websocket.ErrBadHandshake {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return fmt.Errorf("dial failed with status: %d, body: %v", resp.StatusCode, buf.String())
	}
	if err != nil {
		return err
	}

	s.Conn = c

	return nil
}

// Handshake with server
func (s *Socket) Handshake() error {
	type Payload struct {
		Type    string `json:"type"`
		Version string `json:"version"`
		Key     string `json:"key"`
	}
	type Handshake struct {
		Method  string   `json:"method"`
		Payload *Payload `json:"payload"`
	}

	hs := &Handshake{
		Method: "handshake",
		Payload: &Payload{
			Type:    "client",
			Version: "2.0",
			Key:     "17paIsICur8sA0OBqG6dH5G1rmrHNMwt4oNk4iX9",
		},
	}
	h, err := json.Marshal(hs)
	if err != nil {
		return err
	}

	if err := s.Conn.WriteMessage(websocket.TextMessage, h); err != nil {
		return err
	}

	return nil
}

// Monitor the websocket for messages
func (s *Socket) Monitor(done chan struct{}, response chan SocketResponse) {
	go func() {
		defer close(done)

		for {
			sr := SocketResponse{}

			_, message, err := s.Conn.ReadMessage()
			if err != nil {
				sr.Error = err
				response <- sr
				continue
			}

			m := &Method{}
			if err := json.Unmarshal(message, m); err != nil {
				sr.Error = err
				response <- sr
				continue
			}

			switch m.Method {
			case "handshake":
				log.Println("Handshake successful")
			case "notifyTradesInserted":
				p := &struct {
					Payload []*TradeInserted `json:"payload"`
				}{}
				if err := json.Unmarshal(message, p); err != nil {
					sr.Error = fmt.Errorf("unmarshal error: %v, for message %v", err, string(message))
					response <- sr
				} else {
					for _, t := range p.Payload {
						sr.TradeInserted = t
						response <- sr
					}
				}
			case "notifyOrderInserted":
				p := &struct {
					Payload *OrderInserted `json:"payload"`
				}{}
				if err := json.Unmarshal(message, p); err != nil {
					sr.Error = fmt.Errorf("unmarshal error: %v, for message %v", err, string(message))
				} else {
					sr.OrderInserted = p.Payload
				}
				response <- sr
			case "pushCancel":
				p := &struct {
					Payload *PushCancel `json:"payload"`
				}{}
				if err := json.Unmarshal(message, p); err != nil {
					sr.Error = fmt.Errorf("unmarshal error: %v, for message %v", err, string(message))
				} else {
					sr.PushCancel = p.Payload
				}
				response <- sr
			case "pushCancels":
				p := &struct {
					Payload []*PushCancel `json:"payload"`
				}{}
				if err := json.Unmarshal(message, p); err != nil {
					sr.Error = fmt.Errorf("unmarshal error: %v, for message %v", err, string(message))
					response <- sr
				} else {
					for _, pc := range p.Payload {
						sr.PushCancel = pc
						response <- sr
					}
				}
			case "pushEthPrice":
			case "pushServerBlock":
			case "pushRewardPoolSize":
			default:
				log.Printf("Other method: %+v\n", string(message))
			}

		}
	}()
}
