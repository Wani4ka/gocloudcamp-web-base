package main

import (
	"gocloudcamp/proto"
	"gocloudcamp/server"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8089"
	}
	log.Printf("Listening on port %v", port)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("cannot create listener: %v", err)
	}
	serverRegistar := grpc.NewServer()
	proto.RegisterPlaylistServer(serverRegistar, server.NewPlaylistServer())
	err = serverRegistar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %v", err)
	}
}
