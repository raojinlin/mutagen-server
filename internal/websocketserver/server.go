package websocketserver

import (
	"github.com/gorilla/websocket"
	"github.com/raojinlin/mutagen-server/internal/logger"
	"os"
	"reflect"
	"sync"
)

type Request struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
	Id     int         `json:"id"`
}

type Response struct {
	Message string      `json:"message"`
	Action  string      `json:"action"`
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Id      int         `json:"id"`
}

type Handler func(request Request) (interface{}, *Error)

// Server - empty struct for server
type Server struct {
	Conn     *websocket.Conn
	handlers map[string][]Handler
	logger   *logger.Logger
	sync.Mutex
}

func (s *Server) Message(message Response) {
	s.Lock()
	err := s.Conn.WriteJSON(message)
	if err != nil {
		s.logger.Error().Println("write JSON error", err)
	}
	s.Unlock()
}

func (s *Server) Handle(action string, handler Handler) {
	s.Register(action, handler)
}

func (s *Server) Register(action string, handler Handler) {
	if _, ok := s.handlers[action]; !ok {
		s.handlers[action] = []Handler{}
	}

	s.logger.Println("Register action handler: ", action)
	s.handlers[action] = append(s.handlers[action], handler)
}

func (s *Server) dispatch(req Request) {
	if handlers, ok := s.handlers[req.Action]; ok {
		wg := sync.WaitGroup{}
		wg.Add(len(handlers))
		for _, handler := range handlers {
			go func(handler2 Handler) {
				defer wg.Done()
				defer func() {
					err := recover()
					if err == nil {
						return
					}

					str := reflect.ValueOf(err).String()
					s.Message(Response{
						Message: str,
						Action:  "error",
						Data:    nil,
						Code:    CodeInternalError,
						Id:      0,
					})
					s.logger.Error().Println(err)
				}()
				data, err := handler2(req)
				if err != nil {
					s.Message(Response{
						Message: err.Error(),
						Action:  "error",
						Data:    nil,
						Code:    err.Code,
						Id:      req.Id,
					})
				} else {
					s.Message(Response{
						Action: "ack",
						Data:   data,
						Id:     req.Id,
					})
				}
			}(handler)
		}

		wg.Wait()
	}
}

func (s *Server) Start() error {
	for {
		var p Request
		err := s.Conn.ReadJSON(&p)
		s.logger.Println("read json => ", p)
		if err == nil {
			go func() {
				s.dispatch(p)
			}()
		} else {
			return err
		}
	}
}

func NewServer(conn *websocket.Conn) *Server {
	return &Server{
		Conn:     conn,
		handlers: map[string][]Handler{},
		logger:   logger.NewLogger("Websocket-Server", os.Stdout, logger.LevelInfo, false),
	}
}
