package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	hello "grpc-study/pkg/grpc"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

type MyServiceServer struct {
	hello.UnimplementedGreetingServiceServer
}

func (s *MyServiceServer) Hello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	return &hello.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func (s *MyServiceServer) HelloServerStream(req *hello.HelloRequest, stream hello.GreetingService_HelloServerStreamServer) error {
	resCount := 5
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&hello.HelloResponse{
			Message: fmt.Sprintf("[%d] Hello, %s!", i, req.GetName()),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}
func (s *MyServiceServer) HelloClientStream(stream hello.GreetingService_HelloClientStreamServer) error {
	list := make([]string, 0)

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Printf("end of stream")
			break
		}
		if err != nil {
			return err
		}
		log.Printf("request: %v", req)
		list = append(list, req.GetName())
	}
	return stream.SendAndClose(&hello.HelloResponse{
		Message: fmt.Sprintf("Hello, %v!", list),
	})
}

func (s *MyServiceServer) HelloBiDiStream(stream hello.GreetingService_HelloBiDiStreamServer) error {
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Printf("end of stream")
			break
		}
		if err != nil {
			return err
		}
		log.Printf("request: %v", req)
		if err := stream.Send(&hello.HelloResponse{
			Message: fmt.Sprintf("Hello, %s!", req.GetName()),
		}); err != nil {
			return err
		}
	}
	return nil
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
