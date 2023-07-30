package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	hello "grpc-study/pkg/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

type MyServiceServer struct {
	hello.UnimplementedGreetingServiceServer
}

func (s *MyServiceServer) Hello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	return &hello.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func NewServiceServer() *MyServiceServer {
	return &MyServiceServer{}
}

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	hello.RegisterGreetingServiceServer(s, NewServiceServer())

	reflection.Register(s)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
