package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitRoutes(grpcConn *grpc.ClientConn, c *gin.Engine) {
	InitSyncRoutes(grpcConn, c)
	// index
	c.GET("/", func(c *gin.Context) {
		c.Writer.Write([]byte("home"))
	})
}
