package main

import (
	"gocloudcamp/core/playlist"
	"gocloudcamp/proto"
	"gocloudcamp/server/crud"
	"gocloudcamp/server/seek"
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
	stored := playlist.NewPlaylist()
	proto.RegisterCRUDServer(serverRegistar, crud.NewServer(stored))
	proto.RegisterSeekServer(serverRegistar, seek.NewServer(stored))
	err = serverRegistar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %v", err)
	}
}
