class Client {
    constructor(url) {
        const websocket = new WebSocket("ws://127.0.0.1:8081/synchronization")
        websocket.onopen = function () {
            console.log(url, "connected");
        };

        websocket.onmessage = (event)  => {
            let { data } = event;
            data = JSON.parse(data)
            if (data.action === "prompt") {
                console.log(data.message);
                this.sendMessage("promptAck", prompt(data.message))
            } else if (data.action === "message") {
                console.log(data.message);
            } else if (data.action === "error") {
                console.error(data.message)
            } else {
                console.log(data)
            }
        }
        websocket.onclose = function () {
            console.log("CLOSED")
        }
        this.websocket = websocket
        this.id = 1
    }

    sendMessage(action, data) {
        this.websocket.send(JSON.stringify({
            action,
            id: this.id++,
            data
        }))
    }
}