# Mutagen Server

这是一个为```mutagen```提供websocket接口的服务器软件。目前只支持文件同步接口，支持SSH密码输入。

## 接口
* 文件同步
  * 查询
  * 创建
  * 暂停
  * 恢复
  * 重置
  * 终结

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

websocket接口请求与响应数据格式如下。

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
创建参数，参考[synchronization.pb.go:27](github.com/mutagen-io/mutagen@v0.16.2/pkg/service/synchronization/synchronization.pb.go:27)
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