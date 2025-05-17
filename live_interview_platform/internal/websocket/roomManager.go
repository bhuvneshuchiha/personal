package websocket

import "sync"

type RoomManager struct {
	Rooms map[string]*Room
	Mu sync.Mutex
}

func (r *RoomManager) CreateRoom() {
}


func (r *RoomManager) GetRoom() {
}


func (r *RoomManager) DeleteRoom() {
}
