package api

import (
	"github.com/gin-contrib/cors"
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

func initRoutes(grpcConn *grpc.ClientConn, c *gin.Engine) {
	initSyncRoutes(grpcConn, c)
	// index
	c.GET("/", func(c *gin.Context) {
		c.Writer.Write([]byte("home"))
	})
}

func SetupRouter(cc *grpc.ClientConn, listen string) *gin.Engine {
	router := gin.Default()

	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = []string{"http://" + listen, "http://localhost:3000"}

	router.Use(cors.New(cfg))

	initRoutes(cc, router)
	return router
}
