package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ngoctd314/grpc-learning/features/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serviceConfig = `{
	"loadBalancingPolicy": "round_robin",
	"healthCheckConfig": {
		"serviceName": ""
	}
}`

func callUnaryEcho(c pb.EchoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.UnaryEcho(ctx, &pb.EchoRequest{
		Message: fmt.Sprintf("now: %d", time.Now().UnixNano()),
	})
	if err != nil {
		fmt.Println("UnaryEcho: _", err, r)
	} else {
		fmt.Println("UnaryEcho: ", r.GetMessage())
	}
}

func main() {
	flag.Parse()

	address := fmt.Sprintf("localhost:8081")

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}
	defer conn.Close()

	echoClient := pb.NewEchoClient(conn)

	for {
		callUnaryEcho(echoClient)
		time.Sleep(time.Second)
	}
}
