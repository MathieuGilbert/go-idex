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
func (s *Socket) Monitor(done chan struct{}, resp chan SocketResponse) {
	defer close(done)

	for {
		sr := SocketResponse{}

		_, msg, err := s.Conn.ReadMessage()
		if err != nil {
			sr.Error = err
			resp <- sr
			continue
		}

		m := &Method{}
		if err := json.Unmarshal(msg, m); err != nil {
			sr.Error = err
			resp <- sr
			continue
		}

		switch m.Method {
		case "handshake":
			log.Println("Handshake successful")
		case "notifyTradesInserted":
			processTradesInserted(msg, resp)
		case "notifyOrderInserted":
			processOrderInserted(msg, resp)
		case "pushCancel":
			processPushCancel(msg, resp)
		case "pushCancels":
			processPushCancels(msg, resp)
		case "pushEthPrice":
		case "pushServerBlock":
		case "pushRewardPoolSize":
		default:
			log.Printf("Other method: %+v\n", string(msg))
		}

	}
}

func processTradesInserted(msg []byte, c chan SocketResponse) {
	p := &struct {
		Payload []*TradeInserted `json:"payload"`
	}{}
	sr := SocketResponse{}

	if err := json.Unmarshal(msg, p); err != nil {
		sr.Error = fmt.Errorf("unmarshal error: %v, for message %v", err, string(msg))
		c <- sr
	} else {
		for _, t := range p.Payload {
			sr.TradeInserted = t
			c <- sr
		}
	}
}

func processOrderInserted(msg []byte, c chan SocketResponse) {
	p := &struct {
		Payload *OrderInserted `json:"payload"`
	}{}
	sr := SocketResponse{}

	if err := json.Unmarshal(msg, p); err != nil {
		sr.Error = fmt.Errorf("unmarshal error: %v, for message %v", err, string(msg))
	} else {
		sr.OrderInserted = p.Payload
	}
	c <- sr
}

func processPushCancel(msg []byte, c chan SocketResponse) {
	p := &struct {
		Payload *PushCancel `json:"payload"`
	}{}
	sr := SocketResponse{}

	if err := json.Unmarshal(msg, p); err != nil {
		sr.Error = fmt.Errorf("unmarshal error: %v, for message %v", err, string(msg))
	} else {
		sr.PushCancel = p.Payload
	}
	c <- sr
}

func processPushCancels(msg []byte, c chan SocketResponse) {
	p := &struct {
		Payload []*PushCancel `json:"payload"`
	}{}
	sr := SocketResponse{}

	if err := json.Unmarshal(msg, p); err != nil {
		sr.Error = fmt.Errorf("unmarshal error: %v, for message %v", err, string(msg))
		c <- sr
	} else {
		for _, pc := range p.Payload {
			sr.PushCancel = pc
			c <- sr
		}
	}
}
