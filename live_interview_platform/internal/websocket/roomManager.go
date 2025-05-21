package websocket

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

type RoomManager struct {
	Rooms map[string]*Room
	Mu    sync.Mutex
}

func (r *RoomManager) CreateRoom() string {
	fmt.Println("Inside the create room function")
	roomInstance := Room{
		ID:         uuid.New(),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		CodeState:  "",
		Mu:         sync.Mutex{},
	}
	//Add the roomInstance to the roomManager
	r.Mu.Lock()
	r.Rooms[roomInstance.ID.String()] = &roomInstance
	r.Mu.Unlock()

	go func() {
		for {
			select {
			case val := <-roomInstance.Register:
				roomInstance.Mu.Lock()
				roomInstance.Clients[val] = true
				roomInstance.Mu.Unlock()
				log.Printf("New client registered: %s in room %s", val.ID, roomInstance.ID)

			case val := <-roomInstance.Unregister:
				roomInstance.Mu.Lock()
				delete(roomInstance.Clients, val)
				roomInstance.Mu.Unlock()

			case msg := <-roomInstance.Broadcast:
				for client := range roomInstance.Clients {
					select {
					case client.Send <- msg:
					default:
						// if we can't send, remove the client
						roomInstance.Mu.Lock()
						delete(roomInstance.Clients, client)
						roomInstance.Mu.Unlock()
						close(client.Send)
					}
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
