package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/thewebdevel/grpc-file-transfer/messaging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (*server) Upload(stream messaging.FileTransferService_UploadServer) error {
	return nil
}

func main() {
	// If we get the crash code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Start Service
	fmt.Println("File Transfer service started")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	messaging.RegisterFileTransferServiceServer(s, &server{})

	// Register a reflection
	reflection.Register(s)

	go func() {
		fmt.Println("Starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for ctrl c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until the signal is received
	<-ch
	fmt.Println("Stopping the server...")
	s.Stop()
	fmt.Println("Closing the listener...")
	lis.Close()
	fmt.Println("End of Program")
}
