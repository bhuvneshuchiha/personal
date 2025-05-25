package client

import (
	"github.com/bhuvneshuhciha/project_mordoria/internal/message"
	"github.com/google/uuid"
)

type Client struct {
	ID uuid.UUID
	SendMessage chan *message.Message
	MessagesCount int
}

