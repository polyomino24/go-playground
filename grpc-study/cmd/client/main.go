package main

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	hello "grpc-study/pkg/grpc"
	"io"
	"log"
	"os"
)

var (
	scanner *bufio.Scanner
	client  hello.GreetingServiceClient
)

func main() {
	fmt.Println("start gRPC Client.")

	scanner = bufio.NewScanner(os.Stdin)

	address := "localhost:8080"
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("Connection failed.")
		return
	}
	defer conn.Close()

	client = hello.NewGreetingServiceClient(conn)

	for {
		fmt.Println("1: send Request")
		fmt.Println("2: HelloServerStream")
		fmt.Println("3: HelloClientStream")
		fmt.Println("4: HelloBiDiStream")
		fmt.Println("5: exit")
		fmt.Print("please enter >")

		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			Hello()

		case "2":
			HelloServerStream()

		case "3":
			HelloClientStream()
		case "4":
			HelloBiDiStream()

		case "5":
			fmt.Println("bye.")
			goto M
		}
	}
M:
}

func Hello() {
	fmt.Println("Please enter your name.")
	scanner.Scan()
	name := scanner.Text()

	req := &hello.HelloRequest{
		Name: name,
	}
	res, err := client.Hello(context.Background(), req)
	if err != nil {
		if a, ok := status.FromError(err); ok {
			fmt.Println(a.Details())
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(res.GetMessage())
	}
}

func HelloServerStream() {
	fmt.Println("Please enter your name.")
	scanner.Scan()
	name := scanner.Text()

	req := &hello.HelloRequest{
		Name: name,
	}
	stream, err := client.HelloServerStream(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("all the responses have already received.")
			break
		}

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
}

func HelloClientStream() {
	stream, err := client.HelloClientStream(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Please enter your name.")
	for {
		scanner.Scan()
		name := scanner.Text()
		if name == "exit" {
			break
		}
		if err := stream.Send(&hello.HelloRequest{
			Name: name,
		}); err != nil {
			fmt.Println(err)
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}

func HelloBiDiStream() {
	stream, err := client.HelloBiDiStream(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Please enter your name.")
	for {
		scanner.Scan()
		name := scanner.Text()
		if name == "exit" {
			break
		}
		if err := stream.Send(&hello.HelloRequest{
			Name: name,
		}); err != nil {
			fmt.Println("err")
			return
		}
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("err1")
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)
	}

	if err := stream.CloseSend(); err != nil {
		fmt.Println(err)
		return
	}
}
