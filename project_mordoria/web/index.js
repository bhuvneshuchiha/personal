const WebSocket = require('ws');

const ws = new WebSocket('ws://localhost:8080/ws/v1/mordoria');

ws.on('open', () => {
    console.log('Connected to Go backend');
    ws.send(JSON.stringify({
        clientId: 1,
        payload: 'Hey from javascript'
    }));
});

ws.on('message', (data) => {
    console.log('Received from Go backend:', data.toString());
});

ws.on('close', () => {
    console.log('Disconnected from Go backend');
});

ws.on('error', (error) => {
    console.error('WebSocket error:', error);
});

