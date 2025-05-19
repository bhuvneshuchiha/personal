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
			case val :=  <- roomInstance.Register:
				roomInstance.Clients[val] = true
			case val := <- roomInstance.Unregister:
				delete(roomInstance.Clients, val)
			case msg :=  <- roomInstance.Broadcast:
				for client := range roomInstance.Clients {
					client.Send <- msg
				}
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
	roomInst := r.Rooms[roomId]
	r.Mu.Unlock()

	if roomInst != nil {
		go func() {
			roomInst.Register <- client
		}()
		return "Client added to the room"
	}
	return "Room Id does not exist"
}

func (r *RoomManager) UnregisterClient(roomId string, client *Client) string {

	r.Mu.Lock()
	roomInst := r.Rooms[roomId]
	r.Mu.Unlock()

	if roomInst != nil {
		go func() {
			roomInst.Unregister <- client
		}()
		return "Successfully unregistered the client"
	}
	return "Room ID does not exist"
}

func (r *RoomManager) BroadcastToRoom(roomId string, message *Message) string {

	r.Mu.Lock()
	roomInst := r.Rooms[roomId]
	r.Mu.Unlock()

	if roomInst != nil {
		go func() {
			roomInst.Broadcast <- message
		}()
		return "Message added to broadcast queue."
	}
	return "Room Id does not exist"
}


