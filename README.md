# Mutagen Server

该项目是一个为Mutagen提供WebSocket接口的服务器软件。目前，该项目只支持文件同步接口，并支持SSH密码输入。该项目旨在为Mutagen做一个Web端扩展，允许用户通过GUI创建和监控同步任务。

## 安装
使用以下命令进行安装：

```bash
$ go get -u github.com/raojinlin/mutagen-server
```

## 使用示例

命令参数
```bash
$ mutagen-server --help
Usage of mutagen-server
  -listen string
        specify listen address (default "127.0.0.1:8081")
  -sock string
        specify daemon sock path (default "unix:/Users/raojinlin/.mutagen/daemon/daemon.sock")

```

启动服务器
```bash
$ mutagen-server -listen 127.0.0.1:8081 -sock ~/.mutagen/daemon//daemon.sock
```
你可以使用curl命令测试是否启动成功。
```bash
$ curl 127.0.0.1:8081
```

## 支持的接口
* 文件同步
  * 查询同步会话
  * 创建同步任务
  * 暂停同步任务
  * 恢复同步任务
  * 重置同步任务
  * 终止同步任务

### 查询同步会话
```http request
GET /api/synchronization/sessions?id=SESSION_ID&label=LABEL_SELECTOR&name=NAME
```

#### 参数
* SESSION_ID: 会话ID
* LABEL_SELECTOR: 标签选择器，```LABEL=VALUE```，与```SESSION_ID```和```NAME```参数互斥，优先取```SESSION_ID```和```NAME```

#### 会话列表响应
```json
{
    "message": "",
    "action": "list",
    "data": {
        "stateIndex": 427168,
        "sessionStates": [
            {
                "session": {
                    "identifier": "sync_cTNbmVcnIfQ5jrSfXLCZdW0UDIu7zWOfhH9Pu6Bzs1H",
                    "version": 1,
                    "creationTime": {
                        "seconds": 1670856053,
                        "nanos": 168591000
                    },
                    "creatingVersionMinor": 15,
                    "creatingVersionPatch": 2,
                    "alpha": {
                        "path": "/tmp/mugaten-test1"
                    },
                    "beta": {
                        "protocol": "ssh",
                        "user": "raojinlin",
                        "host": "172.31.2.19",
                        "path": "/tmp/mutagen-sync"
                    },
                    "configuration": {},
                    "configurationAlpha": {},
                    "configurationBeta": {},
                    "name": "test",
                    "labels": {
                        "env": "test",
                        "for": "ui",
                        "name": "test"
                    }
                },
                "status": "connecting-beta",
                "lastError": "beta scan error: unable to receive scan response: unable to read message length: unexpected EOF",
                "alphaState": {
                    "connected": true
                },
                "betaState": {}
            }
        ]
    },
    "code": 0,
    "id": 0
}
```

### Websocket接口

下面是文件同步接口配置，websocket接口请求与响应数据格式如下。

请求
```go
type Request struct {
	// 接口名称
	Action  string      `json:"action"`
	// 接口数据
	Data    interface{} `json:"data"`
	// 请求ID
	Id      int         `json:"id"`
}
```
响应
```go
type Response struct {
	// 接口名称
	Action  string      `json:"action"`
	// 如果错误码大于0则表示有错误发生，这是错误信息	
	Message string      `json:"message"`
	// 正常情况下会接口的返回值，当有错误发生时为null   
	Data    interface{} `json:"data"`
	// 错误码，0表示正常，没有错误
	Code    int         `json:"code"`
	// 请求ID，用来标识请求响应ID与请求ID一致
	Id      int         `json:"id"`
}
```


#### 接口概览表

| Action    | 描述          |
|-----------|-------------|
| creation  | 创建同步        |
| list      | 查询已经创建的同步会话 |
| pause     | 暂停会话        |
| reset     | 重置会话        |
| resume    | 恢复暂停的会话     |
| terminate | 终止会话        |


#### 创建同步

请求
```json
{
  "action": "creation",
  "id": 1,
  "data": {}
}
```

响应
创建参数，参考[synchronization.pb.go:27](https://github.com/mutagen-io/mutagen/tree/master/pkg/service/synchronization/synchronization.pb.go#L27)
```json
{
  "action": "creation",
  "id": 1,
  "data": {
    "name": "",
    "labels": [],
    "paused": false,
    "alpha": {
      "kind": 0,
      "protocol": 0,
      "user": "",
      "port": 0,
      "environment": {},
      "parameters": {},
      "path": ""
    },
    "beta": {
      "kind": 0,
      "protocol": 0,
      "user": "",
      "port": 0,
      "environment": {},
      "parameters": {},
      "path": ""
    },
    "configuration": {
      "synchronizationMode":    0,
      "maximumEntryCount":      0,
      "maximumStagingFileSize": 0,
      "probeMode":              0,
      "scanMode":               0,
      "stageMode":              0,
      "symbolicLinkMode":       0,
      "watchMode":              0,
      "watchPollingInterval":   0,
      "defaultIgnores":         [],
      "ignores":                [],
      "ignoreVCSMode":          false,
      "permissionsMode":        0,
      "defaultFileMode":        0,
      "defaultDirectoryMode":   0,
      "defaultOwner":           "",
      "defaultGroup":           ""
    },
    "alphaConfiguration": {
      "synchronizationMode":    0,
      "maximumEntryCount":      0,
      "maximumStagingFileSize": 0,
      "probeMode":              0,
      "scanMode":               0,
      "stageMode":              0,
      "symbolicLinkMode":       0,
      "watchMode":              0,
      "watchPollingInterval":   0,
      "defaultIgnores":         [],
      "ignores":                [],
      "ignoreVCSMode":          false,
      "permissionsMode":        0,
      "defaultFileMode":        0,
      "defaultDirectoryMode":   0,
      "defaultOwner":           "",
      "defaultGroup":           ""
    },
    "betaConfiguration": {
      "synchronizationMode":    0,
      "maximumEntryCount":      0,
      "maximumStagingFileSize": 0,
      "probeMode":              0,
      "scanMode":               0,
      "stageMode":              0,
      "symbolicLinkMode":       0,
      "watchMode":              0,
      "watchPollingInterval":   0,
      "defaultIgnores":         [],
      "ignores":                [],
      "ignoreVCSMode":          0,
      "permissionsMode":        0,
      "defaultFileMode":        0,
      "defaultDirectoryMode":   0,
      "defaultOwner":           "",
      "defaultGroup":           ""
    }
  },
  "error": 0,
  "message": ""
}
```

#### 查询
```json
{
  "action": "list",
  "data": {
    "name": "test",
    "id": "xxx",
    "label": "x=1"
  }
}
```

#### 暂停
```json
{
  "action": "pause",
  "data": {
    "name": "test",
    "id": "xxx",
    "label": "x=1"
  }
}
```

#### 停止
```json
{
  "action": "terminate",
  "data": {
    "name": "test",
    "id": "xxx",
    "label": "x=1"
  }
}
```


#### 恢复
```json
{
  "action": "resume",
  "data": {
    "name": "test",
    "id": "xxx",
    "label": "x=1"
  }
}
```


#### 重置
```json
{
  "action": "reset",
  "data": {
    "name": "test",
    "id": "xxx",
    "label": "x=1"
  }
}
```

### 客户端示例

```javascript
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
```