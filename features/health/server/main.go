package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ngoctd314/grpc-learning/features/pb"
	"google.golang.org/grpc"
)

/*
gRPC provides a health library to communicate a system's to their clients. It works by providing a service definition via the health/v1 api

Servers control their serving status. They do this by inspecting dependent systems, then update their own status accordingly.  UNKNOWN,
SERVING, NOT_SERVING and SERVICE_UNKNOWN
*/

var (
	port   = flag.Int("port", 8080, "the port to serve on")
	sleep  = flag.Duration("sleep", time.Second*5, "duration between changes in health")
	system = ""
)

type echoServer struct {
	pb.UnimplementedEchoServer
}

func (s *echoServer) UnaryEcho(ctx context.Context, request *pb.EchoRequest) (*pb.EchoResponse, error) {
	log.Println("RUN")
	msg := request.GetMessage()
	return &pb.EchoResponse{
		Message: fmt.Sprintf("receive: %s, now: %d", msg, time.Now().UnixNano()),
	}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8081))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	// healthcheck := health.NewServer()
	// healthgrpc.RegisterHealthServer(s, healthcheck)

	pb.RegisterEchoServer(grpcServer, &echoServer{})

	// go func() {
	// 	// asynchronously inspect dependencies
	// 	next := healthpb.HealthCheckResponse_SERVING

	// 	for {
	// 		healthcheck.SetServingStatus(system, next)

	// 		if next == healthpb.HealthCheckResponse_SERVING {
	// 			next = healthpb.HealthCheckResponse_NOT_SERVING
	// 		} else {
	// 			next = healthpb.HealthCheckResponse_SERVING
	// 		}

	// 		time.Sleep(*sleep)
	// 	}
	// }()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
