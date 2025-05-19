package websocket

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)


type RoomManager struct {
	Rooms map[string]*Room
	Mu sync.Mutex
}


func (r *RoomManager) CreateRoom() string {
	roomInstance := Room {
		ID: uuid.New(),
		Clients: make(map[*Client]bool),
		Broadcast: make(chan *Message),
		Register: make(chan *Client),
		Unregister: make(chan *Client),
		CodeState: "",
		Mu: sync.Mutex{},
	}
	//Add the roomInstance to the roomManager
	r.Mu.Lock()
	r.Rooms[roomInstance.ID.String()] = &roomInstance
	r.Mu.Unlock()

	go func() {
		for {
			select {
			case <- roomInstance.Register:
				fmt.Println("Someone added to the room")
			case <- roomInstance.Unregister:
				fmt.Println("Someone left the room")
			case val :=  <- roomInstance.Broadcast:
				fmt.Println("Message received", val)
			}
		}
	}()

	return roomInstance.ID.String()
}

func (r *RoomManager) DeleteRoom(roomId string) bool {

	r.Mu.Lock()
	defer r.Mu.Unlock()

	if r.Rooms[roomId] != nil {
		delete(r.Rooms, roomId)
		return true
	}
	return false
}

func (r *RoomManager) RegisterClient(roomId string, client *Client) string {

	r.Mu.Lock()
	defer r.Mu.Unlock()

	if roomId != "" {
		roomInst := r.Rooms[roomId]
		roomInst.Clients[client] = true
		return "Client added to the room"
	}
	return "Room Id does not exist"
}

func (r *RoomManager) UnregisterClient(roomId string, client *Client) string {

	r.Mu.Lock()
	defer r.Mu.Unlock()

	if roomId != "" {
		roomInst := r.Rooms[roomId]
		if roomInst != nil {
			delete(roomInst.Clients, client)
			return "Successfully unregistered the client"
		}
	}
	return "Room ID does not exist"
}

func (r *RoomManager) BroadcastToRoom(roomId string, message *Message) {
}














