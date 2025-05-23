
1. Interviewer: Starts the Session
Visits the app

Lands on a dashboard or landing page

Logs in as "Interviewer"

Enters email/password or uses a basic auth method

Receives a JWT/session token

Creates a new interview room

Clicks “Start Interview”

Server creates a unique sessionID (UUID) and stores session metadata

Interviewer is redirected to:
https://platform.com/interview/{sessionID}

Shares session link with the candidate

Sends the room URL to the candidate manually (email/slack/etc.)

🟨 2. Candidate: Joins the Session
Receives session link

Clicks: https://platform.com/interview/{sessionID}

Logs in as "Candidate"

Enters name or email (no need for complex auth here)

Gets connected to the same room via WebSocket

🧠 3. Interview Room: Core Features
Now both are in the same interview session:

Feature	What Happens
Real-time code editor	Both see and edit the same code buffer, synced over WebSocket
Real-time chat	They can chat; messages are broadcast via WebSocket + stored in DB
Code execution	Candidate or interviewer can click "Run Code" — Go backend runs it
Roles	Interviewer can observe only or edit code; candidate has limited rights
User status	Join/leave messages appear in chat or UI

💾 4. Session Logging

When the session ends:
All chat messages are stored
Final code buffer is saved
Timestamps, user roles, and run history are recorded in DB
Metadata like session duration, language used, etc., may be stored

Interviewer can later:
View the session logs
Replay the code + chat history

🧪 5. Testing & Reliability
Handles disconnections (e.g., user refreshes)
On reconnect, rejoin WebSocket session
Restore previous state (code, chat)

📌 Optional Advanced Flow (Future):

Feature	Description
Spectator join	Read-only role joins as observer
Session tagging	Tag a session with topics (e.g., “Go basics”)
Candidate feedback	Submit rating or notes after interview
AI Feedback	Summarizes candidate performance

🧩 Summary: Interview Timeline
1. Interviewer logs in → creates session → shares link
2. Candidate joins via link → logs in → enters session
3. Both collaborate via editor + chat
4. Code is run via backend sandbox
5. Session is saved to DB
6. Interviewer can review history later






# . The next steps as of today 22-May-2025

Next Steps Breakdown for the WebSocket Live Interview Project

1. Complete Client Connection Lifecycle
    Handle Client Disconnects Gracefully
    Make sure when a client disconnects (closes the WS connection), you unregister it properly from the room and clean up resources.
    → Test and fix any “use of closed network connection” errors.

    Add Client Ping/Pong (Keepalive)
    Implement WebSocket ping/pong to keep connections alive and detect broken connections faster.

2. Message Handling & Broadcasting
    Prevent Echo Back of Sender’s Messages
    Modify broadcast logic to not send the message back to the sender client, unless you want that.
    This fixes duplicate messages in Postman or other clients.

    Add Message Types & Metadata
    Extend Message struct to support different types (chat, join/leave notifications, code edits, etc.).

    Implement Server-Side Commands or Events
    For example, user joins/leaves, typing indicators, or system messages.

3. Room Management Enhancements
    Add Room Cleanup / Auto Deletion
    Automatically delete rooms when empty to free memory.

    Implement Room Listing API
    Allow clients to query active rooms or their participants.

    Add Room Password or Access Control (Optional)
    Implement simple access control if needed.

4. Client Identification & Authentication
    Assign Client Names or User IDs
    Instead of UUID only, allow clients to provide usernames or authenticate.

    Integrate Authentication (JWT, sessions)
    So only authorized users can create/join rooms.

5. Frontend / Client Improvements
    Build a simple Web UI
    For easier testing, visualization of messages and rooms.

    Improve error handling & reconnect logic
    On the client-side WebSocket connection.

6. Testing and Debugging
    Add unit & integration tests
    For room manager, client registration, message broadcast.

    Stress test with multiple clients
    Ensure stability and concurrency safety.

7. Optional: Persistence & Scaling
    Persist chat history
    Store messages to a database for replay.

    Scaling with multiple server instances
    Use Redis pub/sub or similar to sync rooms across servers.

# . Suggested immediate next step for you right now:

    Fix the duplicate message echo problem by modifying your broadcast so the sender client doesn’t receive their own message twice.

    Implement graceful client unregister on disconnect.

    Add basic ping/pong to WebSocket connections.





# NEXT STEPS AS OF 23-05

🔜 Next Steps
1. Ping/Pong Heartbeat Support (Connection Health)

    Use SetReadDeadline, SetPongHandler on server side

    Client should send periodic ping frames

    Detect dead connections more quickly than relying on write/read errors

2. Room Management Enhancements

    Auto-delete empty rooms when all clients disconnect

    Keep a mapping: roomID → []*Client and clean up when empty

    Allow dynamic room creation/joining from the client side (e.g., via message type "join_room")

3. Client Identity and Roles

    Assign usernames or roles (e.g., interviewer, candidate)

    Attach metadata to clients (can use Client struct)

    Enable the UI to reflect who sent what (name, timestamp, etc.)

4. Command-Based Messaging Protocol

Define a basic protocol over WebSocket messages:

{
  "type": "code_update",
  "content": "func main() {}",
  "sender": "interviewer"
}

Types can include:

    code_update

    chat_message

    join_room

    leave_room

Parse these on the backend and route accordingly.
5. Persistent Code State (Optional Next Stage)

    Use a shared state model for code (in-memory for now, Redis if scaling)

    Apply patching or replace entire code buffer

    Support "code sync" on client (send current full buffer on reconnect)

6. Frontend Integration (Basic HTML/JS)

Build a simple frontend to:

    Join a room

    Type and see real-time updates

    Chat with other users

    Leave and rejoin a session

Use Monaco Editor or CodeMirror for real-time collaborative code input.
7. Security

    Add basic auth/token for each client

    Restrict who can join a room

    Add rate limiting or abuse protection if going public

8. Recording/Playback (Optional)

    Store messages or code state deltas per session

    Enable replay of an interview session

🧭 Suggested Roadmap Summary
Stage	Focus
1	Ping/Pong heartbeat + room cleanup
2	Room & client roles, message types
3	Frontend prototype (basic UI + WebSocket)
4	Real-time code sync (Monaco editor)
5	Identity, auth, roles
6	Optional: persistent state or replay
