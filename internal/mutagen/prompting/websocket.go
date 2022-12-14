package prompting

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/raojinlin/mutagen-server/internal/websocketserver"
	"log"
	"time"
)

type WebsocketPrompter struct {
	Conn       *websocket.Conn
	AnswerChan chan string
	Server     *websocketserver.Server
}

func (w *WebsocketPrompter) sendMessage(action, message string) {
	w.Server.Message(websocketserver.Response{
		Message: message,
		Action:  action,
		Data:    nil,
		Code:    0,
		Id:      0,
	})
}

func (w *WebsocketPrompter) Message(message string) error {
	w.sendMessage("message", message)
	return nil
}

func (w *WebsocketPrompter) Prompt(message string) (string, error) {
	w.sendMessage("prompt", message)
	ctx, cancel := context.WithTimeout(context.TODO(), 106*time.Second)
	defer cancel()
	for {
		select {
		case answer := <-w.AnswerChan:
			log.Println("websocket prompter res: ", answer)
			return answer, nil
		case <-ctx.Done():
			cancel()
			return "", ctx.Err()
		}
	}
}
