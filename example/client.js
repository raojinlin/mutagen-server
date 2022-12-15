class Client {
    constructor(url="ws://127.0.0.1:8081/synchronization") {
        const websocket = new WebSocket(url)
        websocket.onopen = function () {
            console.log(url, "connected");
        };

        this.pendding = {};
        websocket.onmessage = (event)  => {
            let { data } = event;
            data = JSON.parse(data)
            if (this.pendding[data.id]) {
                this.pendding[data.id].resolve(data.data);
            }
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

        const prom = function () {
            let res;
            const p = new Promise(resolve => {
                res = resolve;
            });

            p.resolve = res;
            return p;
        }

        return this.pendding[id] = prom();
    }

    create(creation) {
        return this.sendMessage('creation', creation)
    }

    pause(selections) {
        return this.sendMessage('pause', selections)
    }

    reset(selections) {
        return this.sendMessage('reset', selections)
    }

    flush(selections) {
        return this.sendMessage("flush", selections)
    }

    terminate(selections) {
        return this.sendMessage("terminate", selections)
    }
}

async function test() {
    let creation = {
        "name": "mutagen-server-create1",
        "labels": {
            "created_by": "mutagen-server",
            "owner": "raojinlin",
        },
        "paused": false,
        "alpha": {
            "path": "/tmp/mutagen_test"
        },
        "beta": {
            "protocol": "ssh",
            "user": "raojinlin",
            "host":"192.168.31.111",
            "port": 22,
            "path": "/tmp/xxx_sync"
        }
    };

    const client = new Client();
    client.websocket.onopen = async function () {
        const createResult = client.create(creation)
        console.log(createResult);

        const pausedResult = client.pause({name: 'mutagen-server-create1'})
        console.log(pausedResult);
        console.log(await client.sendMessage('list', {name: 'mutagen-server-create1'}));

        const resumeResult = client.sendMessage('resume', {name: 'mutagen-server-create1'})
        console.log(resumeResult);
        console.log(await client.sendMessage('list', {name: 'mutagen-server-create1'}));

        const resetResult = client.reset({name: 'mutagen-server-create1'})
        console.log(resetResult);
        console.log(await client.sendMessage('list', {name: 'mutagen-server-create1'}));

        const terminateResult = client.terminate({name: 'mutagen-server-create1'})
        console.log(terminateResult)
        console.log(await client.sendMessage('list', {name: 'mutagen-server-create1'}));
    }
}