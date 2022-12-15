package cmd

import (
	"flag"
	"github.com/raojinlin/mutagen-server/api"
	"github.com/raojinlin/mutagen-server/internal/grpc"
)

func Run() {
	listen := "127.0.0.1:8081"
	daemonSock := grpc.DefaultAddress()

	flag.StringVar(&listen, "listen", listen, "specify listen address")
	flag.StringVar(&daemonSock, "sock", daemonSock, "specify daemon sock path")
	flag.Parse()

	cc := grpc.Connect(daemonSock)
	defer func() {
		cc.Close()
	}()

	server := api.SetupRouter(cc, listen)
	if err := server.Run(listen); err != nil {
		panic(err)
	}
}
