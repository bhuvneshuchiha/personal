const socket = new WebSocket("ws://localhost:8080/ws")

socket.addEventListener("open", () => {
    console.log("websocket connection has been opened")
    socket.send("Hello Server")
});

socket.addEventListener("message", (event) => {
    console.log("Message from server -> ", event.data)
});

socket.addEventListener("error", (event) => {
    console.log("Error received -> ", event)
});

socket.addEventListener("closed", () => {
    console.log("The websocket connection is closed")
});

function sendMessage(message) {
    if(socket.readyState == WebSocket.OPEN) {
        socket.send(message)
    }else {
        console.error("websocket is not open")
    }
}
