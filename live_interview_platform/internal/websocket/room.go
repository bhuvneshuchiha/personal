package websocket

import (
	"sync"

	"github.com/google/uuid"
)

type Message struct {
	messageString string
	sender string
}

type Room struct {
	ID uuid.UUID
	Clients map[*Client]bool
	Broadcast chan *Message
	Register chan *Client
	Unregister chan *Client
	CodeState string
	Mu sync.Mutex

}

