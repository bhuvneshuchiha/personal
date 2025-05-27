import { useState } from "react";
import { useEffect } from "react";
import { useRef } from "react";

function ChatComponent() {
	const [clientMessages, setClientMessages] = useState("");
	const [ai_emot, set_ai_emot] = useState("");
	const [messages, setMessages] = useState([]);
	const ws = useRef(null);

	useEffect(() => {
		ws.current = new WebSocket("ws://localhost:8081/ws/v1/mordoria");

		ws.current.onopen = () => {
			console.log("Connected to the go backend");
		};

		ws.current.onmessage = (event) => {
			console.log("Received from the go backend", event.data);
            let messageChat = event.data
			setMessages((prev) => [...prev, messageChat]);
		};

		ws.current.onerror = (error) => {
			console.error("Error encountered", error);
		};

		ws.current.onclose = () => {
			console.log("Websocket connection closed");
		};

		return () => {
			ws.current.close();
		};
	}, []);

	const handleSend = () => {
		if (ws.current && ws.current.readyState === WebSocket.OPEN) {
			ws.current.send(
				JSON.stringify({
					clientId: 1,
					payload: clientMessages,
					ai_emot_score: ai_emot.toString(),
				}),
			);
			setClientMessages("");
		}
	};

	return (
		<div>
			<input
				type="text"
				value={clientMessages}
				onChange={(e) => setClientMessages(e.target.value)}
				placeholder="Type your message"
			/>
			<input
				type="text"
				value={ai_emot}
				onChange={(e) => set_ai_emot(e.target.value)}
				placeholder="AI Emot Score"
			/>
			<button onClick={handleSend}>Send</button>

			<div>
				<h3>Received Messages:</h3>
				<ul>
					{messages.map((msg, index) => (
						<li key={index}>{msg}</li>
					))}
				</ul>
			</div>
		</div>
	);
}


export default ChatComponent
