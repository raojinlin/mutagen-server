package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mutagen-io/mutagen/pkg/selection"
	syncsvc "github.com/mutagen-io/mutagen/pkg/service/synchronization"
	"github.com/mutagen-io/mutagen/pkg/synchronization"
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

func unmarshalRequestData(data interface{}, v interface{}) error {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonStr, &v)
	return err
}

func initSyncRoutes(grpcConn *grpc.ClientConn, c *gin.Engine) {
	// websocket service
	c.GET("/synchronization", func(context *gin.Context) {
		conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			context.AbortWithError(500, err)
			return
		}

		wsr := websocketserver.NewServer(conn)
		promptAnswerChan := make(chan string)
		prompter := &prompting.WebsocketPrompter{
			Conn:       wsr.Conn,
			AnswerChan: promptAnswerChan,
			Server:     wsr,
		}

		wsr.Register(websocketserver.SyncCreation, func(request websocketserver.Request) (interface{}, *websocketserver.Error) {
			var creation syncsvc.CreationSpecification
			if request.Data != nil {
				err := unmarshalRequestData(request.Data, &creation)
				if err != nil {
					return nil, websocketserver.NewError(websocketserver.CodeCommonError, err.Error())
				}
			}

			if creation.Configuration == nil {
				creation.Configuration = &synchronization.Configuration{}
			}

			if creation.ConfigurationAlpha == nil {
				creation.ConfigurationAlpha = creation.Configuration
			}

			if creation.ConfigurationBeta == nil {
				creation.ConfigurationBeta = creation.Configuration
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
				err := unmarshalRequestData(request.Data, &sel)
				if err != nil {
					return nil, websocketserver.NewError(websocketserver.CodeCommonError, err.Error())
				}
			}
			r, err := sync.List(grpcConn, generateSelections(sel.Id, sel.Name, sel.Label))
			if err != nil {
				return nil, websocketserver.NewError(websocketserver.CodeCommonError, err.Error())
			}

			return r, nil
		})

		wsr.Register(websocketserver.PromptAck, func(request websocketserver.Request) (interface{}, *websocketserver.Error) {
			promptAnswerChan <- request.Data.(string)
			return nil, nil
		})

		wsr.Register(websocketserver.SyncResume, func(request websocketserver.Request) (interface{}, *websocketserver.Error) {
			var sel Selections
			if request.Data != nil {
				err := unmarshalRequestData(request.Data, &sel)
				if err != nil {
					return nil, websocketserver.NewError(websocketserver.CodeCommonError, err.Error())
				}
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
				err := unmarshalRequestData(request.Data, &sel)
				if err != nil {
					return nil, websocketserver.NewError(websocketserver.CodeCommonError, err.Error())
				}
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
				err := unmarshalRequestData(request.Data, &sel)
				if err != nil {
					return nil, websocketserver.NewError(websocketserver.CodeCommonError, err.Error())
				}
			}
			result, err := sync.Flush(grpcConn, prompter, generateSelections(sel.Id, sel.Name, sel.Label))

			if err != nil {
				return nil, websocketserver.NewError(websocketserver.CodeCommonError, err.Error())
			}

			return result, nil
		})

		wsr.Register(websocketserver.SyncReset, func(request websocketserver.Request) (interface{}, *websocketserver.Error) {
			var sel Selections
			if request.Data != nil {
				err := unmarshalRequestData(request.Data, &sel)
				if err != nil {
					return nil, websocketserver.NewError(websocketserver.CodeCommonError, err.Error())
				}
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
				err := unmarshalRequestData(request.Data, &sel)
				if err != nil {
					return nil, websocketserver.NewError(websocketserver.CodeCommonError, err.Error())
				}
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
			Action:  "ack",
			Data:    res,
			Code:    0,
			Id:      0,
		})
	})

}
