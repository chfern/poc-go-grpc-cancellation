package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/chfern/poc-go-grpc-cancellation/pong/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type helloService struct {
	proto.UnimplementedHelloServiceServer
}

func (h *helloService) Hello(ctx context.Context, spec *proto.HelloSpec) (*proto.HelloResult, error) {
	time.Sleep(2 * time.Second) // intentionally sleep

	if ctx.Err() != nil { // context should've been canceled at this point
		cause := context.Cause(ctx)
		fmt.Printf("context err: %v\n", ctx.Err())
		fmt.Printf("context err cause: %v\n", cause)

		return nil, errors.New("timeout")
	}

	return &proto.HelloResult{
		Payload: fmt.Sprintf("Hello: %s", spec.Payload),
	}, nil
}

func main() {
	// Start grpc server and register helloService
	s := grpc.NewServer()
	helloService := &helloService{}
	proto.RegisterHelloServiceServer(s, helloService)
	reflection.Register(s)

	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting pong grpc at port 8081")
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
