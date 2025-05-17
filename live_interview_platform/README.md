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
