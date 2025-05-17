package websocket

import (
	"sync"

	"github.com/google/uuid"
)

type Room struct {
	ID uuid.UUID
	Clients map[*Client]bool
	Broadcast chan []byte
	Register chan *Client
	Unregister chan *Client
	CodeState string
	Mu sync.Mutex

}

