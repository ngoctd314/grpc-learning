package main

import (
	"context"
	"io"
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

func printFeatures(stub pb.RouteGuideClient, rect *pb.Rectangle) {
	log.Printf("Looking for features within %v", rect)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := stub.ListFeatures(ctx, rect)
	if err != nil {
		log.Fatalf("client.ListFeatures failed: %v", err)
	}

	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.ListFeatures failed: %v", err)
		}
		log.Printf("Feature: name: %q, point: (%v, %v)", feature.GetName(), feature.GetLocation().GetLongitude(), feature.GetLocation().GetLatitude())
	}
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

	printFeatures(stub, &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	})
}
