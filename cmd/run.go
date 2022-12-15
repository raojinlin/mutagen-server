package cmd

import (
	"github.com/raojinlin/mutagen-server/api"
	"github.com/raojinlin/mutagen-server/internal/grpc"
)

func Run() {
	listen := "127.0.0.1:8081"
	cc := grpc.Connect(grpc.DefaultAddress())
	defer func() {
		cc.Close()
	}()

	server := api.SetupRouter(cc, listen)
	if err := server.Run(listen); err != nil {
		panic(err)
	}
}
