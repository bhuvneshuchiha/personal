 ┌───────────────────────┐
 │   RoomManager (global)│◄─── Keeps track of all rooms
 └───────────────────────┘
             │
             ▼
    ┌────────────────────┐
    │     Room (ID: X)   │◄─── One per interview session
    │--------------------│
    │ - clients map      │◄─── userID → *Client
    │ - broadcast chan   │◄─── receives messages from all clients
    │ - join/leave chan  │
    └────────────────────┘
             ▲
     ┌───────┼────────┐
     │       │        │
     ▼       ▼        ▼
 ┌───────┐ ┌───────┐ ┌───────┐
 │Client1│ │Client2│ │Client3│   (in same room)
 └───────┘ └───────┘ └───────┘
   ▲         ▲         ▲
   │         │         │
 WebSocket WebSocket WebSocket
  Conn      Conn      Conn

