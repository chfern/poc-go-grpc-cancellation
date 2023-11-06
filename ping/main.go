package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/chfern/poc-go-grpc-cancellation/ping/proto"
	pongProto "github.com/chfern/poc-go-grpc-cancellation/pong/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type helloService struct {
	pongClient pongProto.HelloServiceClient
	proto.UnimplementedHelloServiceServer
}

// HelloPong calls pong service and return the result from it
func (h *helloService) HelloPong(ctx context.Context, spec *proto.HelloSpec) (*proto.HelloResult, error) {
	ctxWithDeadline, cancelFn := context.WithTimeoutCause(ctx, 1*time.Second, errors.New("ping timeout exceeded"))
	defer cancelFn()

	pongHelloResult, err := h.pongClient.Hello(ctxWithDeadline, &pongProto.HelloSpec{
		Payload: spec.Payload,
	})
	if err != nil {
		return nil, err
	}

	return &proto.HelloResult{
		Payload: fmt.Sprintf("Result from pong: %s", pongHelloResult.Payload),
	}, nil
}

func main() {
	// Create pong client
	pongConn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	pongClient := pongProto.NewHelloServiceClient(pongConn)

	// Start grpc server and register helloService
	s := grpc.NewServer()
	helloService := &helloService{pongClient: pongClient}
	proto.RegisterHelloServiceServer(s, helloService)
	reflection.Register(s)

	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting ping grpc at port 8080")
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
