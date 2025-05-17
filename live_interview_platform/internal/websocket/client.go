package websocket

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)


type Client struct {
	Conn *websocket.Conn
	Send chan []byte
	Room *Room
	Username string
	ID uuid.UUID
	Role string
}
