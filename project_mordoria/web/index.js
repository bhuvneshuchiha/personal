const ws = new WebSocket('ws://localhost:8081/ws/v1/mordoria');

let messageArr = [];

document.getElementById("sendButton").addEventListener("click" ,() => {
    console.log("Inside the tricker")
    const input = document.getElementById("messageInput")
    const message = input.value;
    ws.send(JSON.stringify({
        clientId: 1,
        payload: message
    }));
    input.value = ""
});

ws.onopen = () => {
    console.log('Connected to Go backend');
};

ws.onmessage = (event) => {
    console.log('Sent to Go backend:', event.data);
    messageArr.push(event.data);
    console.log(messageArr);
};

ws.onerror = (error) => {
    console.error('WebSocket error:', error);
};

ws.onclose = () => {
    console.log('Disconnected from Go backend');
};


// ws.on('open', () => {
//     console.log('Connected to Go backend');
// });
//
// ws.on('message', (data) => {
//     console.log('Received from Go backend:', data.toString());
//     messageArr.push(data.toString())
// });
//
// ws.on('close', () => {
//     console.log('Disconnected from Go backend');
// });
//
// ws.on('error', (error) => {
//     console.error('WebSocket error:', error);
// });
//
