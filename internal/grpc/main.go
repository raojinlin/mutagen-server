package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"path"
)

func DefaultAddress() string {
	homedir, _ := os.UserHomeDir()
	return path.Join("unix://", homedir, ".mutagen/daemon/daemon.sock")
}

func Connect(address string) *grpc.ClientConn {
	if address == "" {
		address = DefaultAddress()
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return conn
}
