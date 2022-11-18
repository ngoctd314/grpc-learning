package main

import (
	"context"
	"log"
	"time"

	"github.com/ngoctd314/grpc-learning/route/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func printFeature(client pb.RouteGuideClient, point *pb.Point) {
	log.Printf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	feature, err := client.GetFeature(ctx, point)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(feature)
}

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	stub := pb.NewRouteGuideClient(conn)
	printFeature(stub, &pb.Point{
		Latitude:  10,
		Longitude: 10,
	})
}
