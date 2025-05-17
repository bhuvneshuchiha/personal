package websocket

import (
	"sync"

	"github.com/google/uuid"
)

var mu sync.Mutex

type RoomManager struct {
	Rooms map[string]*Room
	Mutex sync.Mutex
}


func (r *RoomManager) CreateRoom() string {
	messageCh := make(chan *Message)
	clientCh := make(chan *Client)
	clients := make(map[*Client]bool)

	roomInstance := Room {
		ID: uuid.New(),
		Clients: clients,
		Broadcast: messageCh,
		Register: clientCh,
		Unregister: clientCh,
		CodeState: "",
		Mutex: mu,
	}

	roomInstance.Mutex.Lock()
	r.Rooms[roomInstance.ID.String()] = &roomInstance
	roomInstance.Mutex.Unlock()

	return roomInstance.ID.String()
}


func (r *RoomManager) GetRoom() {
}


func (r *RoomManager) DeleteRoom() {
}
