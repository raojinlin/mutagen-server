package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mutagen-io/mutagen/pkg/selection"
	syncsvc "github.com/mutagen-io/mutagen/pkg/service/synchronization"
	"github.com/raojinlin/mutagen-server/internal/mutagen/prompting"
	"github.com/raojinlin/mutagen-server/internal/mutagen/sync"
	"github.com/raojinlin/mutagen-server/internal/websocketserver"
	"google.golang.org/grpc"
)

func generateSelections(id, name, label string) *selection.Selection {
	result := &selection.Selection{
		All:            false,
		Specifications: []string{},
	}

	if id == "" && name == "" && label == "" {
		result.All = true
		return result
	}

	if id != "" {
		result.Specifications = append(result.Specifications, id)
	}

	if name != "" {
		result.Specifications = append(result.Specifications, name)
	}

	if label != "" && name == "" && id == "" {
		result.LabelSelector = label
	}

	return result
}

func generateSelectionsFromQuery(c *gin.Context) *selection.Selection {
	return generateSelections(c.Query("id"), c.Query("name"), c.Query("label"))
}

type Selections struct {
	Name  string `json:"name"`
	Id    string `json:"id"`
	Label string `json:"label"`
}

func InitSyncRoutes(grpcConn *grpc.ClientConn, c *gin.Engine) {
	// websocket service
	c.GET("/synchronization", func(context *gin.Context) {
		conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			context.AbortWithError(500, err)
			return
		}

		wsr := websocketserver.NewServer(conn)
		promptAckChan := make(chan string)
		prompter := &prompting.WebsocketPrompter{
			Conn:       wsr.Conn,
			AnswerChan: promptAckChan,
			Server:     wsr,
		}

		wsr.Register(websocketserver.SyncCreation, func(request websocketserver.Request) (interface{}, *websocketserver.Error) {
			var creation syncsvc.CreationSpecification
			if request.Data != nil {
				jsonStr, _ := json.Marshal(request.Data)
				err := json.Unmarshal(jsonStr, &creation)
				if err != nil {
					return nil, websocketserver.NewError(websocketserver.CodeCommonError, err.Error())
				}
			}

			r, err := sync.Create(grpcConn, prompter, &creation)
			if err != nil {
				return nil, websocketserver.NewError(websocketserver.CodeRequestDataError, err.Error())
			}

			return r, nil
		})

		wsr.Register(websocketserver.SyncList, func(request websocketserver.Request) (interface{}, *websocketserver.Error) {
			var sel Selections
			if request.Data != nil {
				sel = request.Data.(Selections)
			}
			r, err := sync.List(grpcConn, generateSelections(sel.Id, sel.Name, sel.Label))
			if err != nil {
				return nil, websocketserver.NewError(websocketserver.CodeCommonError, err.Error())
			}

			return r, nil
		})

		wsr.Register(websocketserver.PromptAck, func(request websocketserver.Request) (interface{}, *websocketserver.Error) {
			promptAckChan <- request.Data.(string)
			return nil, nil
		})

		wsr.Register(websocketserver.SyncResume, func(request websocketserver.Request) (interface{}, *websocketserver.Error) {
			var sel Selections
			if request.Data != nil {
				sel = request.Data.(Selections)
			}

			result, err := sync.Resume(grpcConn, prompter, generateSelections(sel.Id, sel.Name, sel.Label))

			if err != nil {
				return nil, websocketserver.NewError(websocketserver.CodeCommonError, err.Error())
			}

			return result, nil
		})

		wsr.Register(websocketserver.SyncPause, func(request websocketserver.Request) (interface{}, *websocketserver.Error) {
			var sel Selections
			if request.Data != nil {
				sel = request.Data.(Selections)
			}
			result, err := sync.Pause(grpcConn, prompter, generateSelections(sel.Id, sel.Name, sel.Label))

			if err != nil {
				return nil, websocketserver.NewError(websocketserver.CodeCommonError, err.Error())
			}

			return result, nil
		})

		wsr.Register(websocketserver.SyncFlush, func(request websocketserver.Request) (interface{}, *websocketserver.Error) {
			var sel Selections
			if request.Data != nil {
				sel = request.Data.(Selections)
			}
			result, err := sync.Flush(grpcConn, prompter, generateSelections(sel.Id, sel.Name, sel.Label))

			if err != nil {
				return nil, websocketserver.NewError(websocketserver.CodeCommonError, err.Error())
			}

			return result, nil
		})

		wsr.Register(websocketserver.SyncReset, func(payload websocketserver.Request) (interface{}, *websocketserver.Error) {
			var sel Selections
			if payload.Data != nil {
				sel = payload.Data.(Selections)
			}
			result, err := sync.Reset(grpcConn, prompter, generateSelections(sel.Id, sel.Name, sel.Label))

			if err != nil {
				return nil, websocketserver.NewError(websocketserver.CodeCommonError, err.Error())
			}

			return result, nil
		})

		wsr.Register(websocketserver.SyncTerminate, func(request websocketserver.Request) (interface{}, *websocketserver.Error) {
			var sel Selections
			if request.Data != nil {
				sel = request.Data.(Selections)
			}
			result, err := sync.Terminate(grpcConn, prompter, generateSelections(sel.Id, sel.Name, sel.Label))

			if err != nil {
				return nil, websocketserver.NewError(websocketserver.CodeCommonError, err.Error())
			}

			return result, nil
		})

		err = wsr.Start()
		if err != nil {
			fmt.Println(err)
			context.Abort()
		}
	})

	// restful api
	c.GET("/api/synchronization/sessions", func(c *gin.Context) {
		res, err := sync.List(grpcConn, generateSelectionsFromQuery(c))

		if err != nil {
			c.IndentedJSON(500, websocketserver.Response{
				Message: err.Error(),
				Action:  "error",
				Data:    nil,
				Code:    websocketserver.CodeCommonError,
				Id:      0,
			})
			return
		}

		c.IndentedJSON(200, websocketserver.Response{
			Message: "",
			Action:  "list",
			Data:    res,
			Code:    0,
			Id:      0,
		})
	})

}
