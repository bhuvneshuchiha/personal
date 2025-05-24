package masterRoom

import (
	"errors"

	"github.com/google/uuid"
)

type MasterRoom struct {
	ID      uuid.UUID          `json:"id"`
	Clients map[string]*Client `json:"clients_table"`
	ListenMessage chan *Message`json:"message_string"`
	RoomTag string             `json:"room_tag"`
}

func (m *MasterRoom) CreateMasterRoom() *MasterRoom {
	return &MasterRoom{
		ID:      uuid.New(),
		Clients: make(map[string]*Client),
		ListenMessage: make(chan *Message),
		RoomTag: "Witty",
	}
}

func (m *MasterRoom) AddClient(client *Client, id string) error {
	if id == "" {
		return errors.New("Invalid ID")
	}

	_, ok := m.Clients[id]
	if !ok {
		return errors.New("Client already attached to room")
	}

	m.Clients[id] = client
	return nil
}


func (m *MasterRoom) RemoveClient(id string) error {
	if id == "" {
		return errors.New("Invalid ID")
	}

	_, ok := m.Clients[id]
	if !ok {
		return errors.New("Client does not exist")
	}

	delete(m.Clients, id)
	return nil
}


func (m *MasterRoom) BroadCastMessage() {

}







