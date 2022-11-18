package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/ngoctd314/grpc-learning/route/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

/*
There are two parts to making our RouteGuide service do its job:
- Implementing the service interface generated from our service definition: doing the actual "work" of our service
- Running a gRPC server to listen for requests from the clients and dispatch them to the right service implementation
*/

type routeGuideServer struct {
	pb.UnimplementedRouteGuideServer
	savedFeatures []*pb.Feature

	mu         sync.Mutex // guards
	routeNotes map[string][]*pb.RouteNote
}

// Then you define rpc methods inside your service definition, specifying
// their request and response types
// A simple RPC where the client sends a
// request to the server using the stub and waits for a response
func (r *routeGuideServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	for _, feature := range r.savedFeatures {
		if proto.Equal(feature.Location, point) {
			return feature, nil
		}
	}

	log.Println("receive: ", point.GetLatitude(), point.GetLongitude())

	// No feature was found, return an unnamed feature
	return &pb.Feature{
		Name: "Test feature",
		Location: &pb.Point{
			Latitude:  15,
			Longitude: 20,
		},
	}, nil
}

// A server-side streaming RPC where the client sends a request to the server
// and gets a stream to read a sequence of message back.
// The client reads from the returned stream until there are no more messages.
// You specify a server-side streaming method by placing the stream keyword
// before the response type
func (r *routeGuideServer) ListFeatures(rect *pb.Rectangle, features pb.RouteGuide_ListFeaturesServer) error {
	panic("not implemented") // TODO: Implement
}

func newServer() *routeGuideServer {
	s := &routeGuideServer{
		routeNotes: map[string][]*pb.RouteNote{},
	}

	return s
}

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterRouteGuideServer(grpcServer, newServer())
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
	if err != nil {
		log.Fatal(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
