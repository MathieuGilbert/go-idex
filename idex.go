package idex

// Idex service
type Idex struct {
	API    *API
	Socket *Socket
}

// IDEX rest and websocket urls
const (
	APIURL = "https://api.idex.market"
	WSURL  = "wss://v1.idex.market"
)

// New instance of an Idex
func New() *Idex {
	return &Idex{API: &API{URL: APIURL}, Socket: &Socket{URL: WSURL}}
}
