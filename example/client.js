class Client {
    constructor(url="ws://127.0.0.1:8081/synchronization") {
        const websocket = new WebSocket(url)
        websocket.onopen = function () {
            console.log(url, "connected");
        };

        websocket.onmessage = (event)  => {
            let { data } = event;
            data = JSON.parse(data)
            if (data.action === "prompt") {
                console.log(data.message);
                this.sendMessage("answer", prompt(data.message))
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
        const id = this.id++;
        this.websocket.send(JSON.stringify({
            action,
            id,
            data
        }))
    }

    create(creation) {
        this.sendMessage('creation', creation)
    }

    pause(selections) {
        this.sendMessage('pause', selections)
    }

    reset(selections) {
        this.sendMessage('reset', selections)
    }

    flush(selections) {
        this.sendMessage("flush", selections)
    }

    terminate(selections) {
        this.sendMessage("terminate", selections)
    }
}