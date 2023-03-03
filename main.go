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
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "storage.local"
	}

	storage := playlist.NewStorage(storagePath, "playlist.dat")
	stored, err := storage.Load()
	if err != nil {
		log.Printf("cannot load a playlist: %v", err)
		log.Println("creating a new one instead")
	}

	log.Printf("listening on port %v", port)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("cannot create listener: %v", err)
	}
	serverRegistar := grpc.NewServer()
	proto.RegisterCRUDServer(serverRegistar, crud.NewServer(stored))
	proto.RegisterSeekServer(serverRegistar, seek.NewServer(stored))
	err = serverRegistar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %v", err)
	}
}
