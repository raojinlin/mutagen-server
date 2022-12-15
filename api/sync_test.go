package api

import (
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"github.com/raojinlin/mutagen-server/internal/grpc"
	"github.com/raojinlin/mutagen-server/internal/websocketserver"
	"net/http/httptest"
	"testing"
)

var cc = grpc.Connect(grpc.DefaultAddress())
var r = SetupRouter(cc, "127.0.0.1:8082")

func TestListAllSyncSessionsApi(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/synchronization/sessions", nil)
	r.ServeHTTP(w, req)

	var response websocketserver.Response
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, "", response.Message)
	assert.Equal(t, "ack", response.Action)
	assert.Equal(t, 0, response.Id)
	assert.NotEqual(t, nil, response.Data)

	data := response.Data.(map[string]interface{})

	assert.Equal(t, 9, len(data["sessionStates"].([]interface{})))
}

func TestListSingleSessionApiWithId(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/synchronization/sessions?id=sync_Vl1bMSyTmhUcf2Xd2O3WrTDpWHQqRBsg68BrR1s9tB7", nil)
	r.ServeHTTP(w, req)

	var response websocketserver.Response
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, "", response.Message)
	assert.Equal(t, "ack", response.Action)
	assert.Equal(t, 0, response.Id)
	assert.NotEqual(t, nil, response.Data)

	data := response.Data.(map[string]interface{})

	states := data["sessionStates"].([]interface{})
	assert.Equal(t, 1, len(states))
	session := states[0].(map[string]interface{})["session"].(map[string]interface{})
	assert.Equal(t, "sync_Vl1bMSyTmhUcf2Xd2O3WrTDpWHQqRBsg68BrR1s9tB7", session["identifier"])
}

func TestListSingleSessionApiWithName(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/synchronization/sessions?name=test", nil)
	r.ServeHTTP(w, req)

	var response websocketserver.Response
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, "", response.Message)
	assert.Equal(t, "ack", response.Action)
	assert.Equal(t, 0, response.Id)
	assert.NotEqual(t, nil, response.Data)

	data := response.Data.(map[string]interface{})

	states := data["sessionStates"].([]interface{})
	assert.Equal(t, 7, len(states))
}

func TestListSingleSessionApiWithLabelEnv(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/synchronization/sessions?label=env=test", nil)
	r.ServeHTTP(w, req)

	var response websocketserver.Response
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, "", response.Message)
	assert.Equal(t, "ack", response.Action)
	assert.Equal(t, 0, response.Id)
	assert.NotEqual(t, nil, response.Data)

	data := response.Data.(map[string]interface{})

	states := data["sessionStates"].([]interface{})
	assert.Equal(t, 7, len(states))
}

func TestListSessionApiWithInvalidName(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/synchronization/sessions?name=xxx", nil)
	r.ServeHTTP(w, req)

	var response websocketserver.Response
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 500, w.Code)
	assert.Equal(t, websocketserver.CodeCommonError, response.Code)
	assert.NotEqual(t, "", response.Message)
	assert.MatchRegex(t, response.Message, "unable to locate requested sessions")
	assert.Equal(t, "error", response.Action)
	assert.Equal(t, 0, response.Id)
	assert.Equal(t, nil, response.Data)
}

func TestListSessionApiWithInvalidLabel(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/synchronization/sessions?label=name=xxx", nil)
	r.ServeHTTP(w, req)

	var response websocketserver.Response
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, "", response.Message)
	assert.Equal(t, "ack", response.Action)
	assert.Equal(t, 0, response.Id)
	assert.NotEqual(t, nil, response.Data)
	//assert.Equal()
	res := response.Data.(map[string]interface{})
	i := res["stateIndex"].(float64)
	assert.Equal(t, true, i > 0)
}
