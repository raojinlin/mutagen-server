package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/raojinlin/mutagen-server/api"
	"github.com/raojinlin/mutagen-server/internal/grpc"
)

func main() {
	server := gin.New()

	cfg := cors.DefaultConfig()
	listen := "127.0.0.1:8081"
	cfg.AllowOrigins = []string{"http://" + listen}

	server.Use(cors.New(cfg))

	cc := grpc.Connect(grpc.DefaultAddress())
	defer func() {
		cc.Close()
	}()

	api.InitRoutes(cc, server)
	if err := server.Run(listen); err != nil {
		panic(err)
	}
}
