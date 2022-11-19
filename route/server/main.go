package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"sync"
	"time"

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
	for _, feature := range r.savedFeatures {
		if inRange(feature.Location, rect) {
			time.Sleep(time.Millisecond * 500)
			if err := features.Send(feature); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *routeGuideServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
	var pointCount, featureCount, distance int32
	var lastPoint *pb.Point
	startTime := time.Now()

	for {
		point, err := stream.Recv()
		if err == io.EOF {
			// close stream
			endTime := time.Now()
			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:   pointCount,
				FeatureCount: featureCount,
				Distance:     distance,
				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}
		pointCount++
		for _, feature := range r.savedFeatures {
			if proto.Equal(feature.Location, point) {
				featureCount++
			}
		}
		if lastPoint != nil {
			distance += calcDistance(lastPoint, point)
		}
		lastPoint = point
	}
}

func (r *routeGuideServer) loadFeatures(filePath string) {
	var data []byte
	if filePath != "" {
		var err error
		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to load default features: %v", err)
		}
	}
	if err := json.Unmarshal(data, &r.savedFeatures); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
}

func inRange(point *pb.Point, rect *pb.Rectangle) bool {
	left := math.Min(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	right := math.Max(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	top := math.Max(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))
	bottom := math.Min(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))

	if float64(point.Longitude) >= left &&
		float64(point.Longitude) <= right &&
		float64(point.Latitude) >= bottom &&
		float64(point.Latitude) <= top {
		return true
	}
	return false
}

func newServer() *routeGuideServer {
	s := &routeGuideServer{
		routeNotes: map[string][]*pb.RouteNote{},
	}
	s.loadFeatures("./db.json")

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

// calcDistance calculates the distance between two points using the "haversine" formula.
// The formula is based on http://mathforum.org/library/drmath/view/51879.html.
func calcDistance(p1 *pb.Point, p2 *pb.Point) int32 {
	const CordFactor float64 = 1e7
	const R = float64(6371000) // earth radius in metres
	lat1 := toRadians(float64(p1.Latitude) / CordFactor)
	lat2 := toRadians(float64(p2.Latitude) / CordFactor)
	lng1 := toRadians(float64(p1.Longitude) / CordFactor)
	lng2 := toRadians(float64(p2.Longitude) / CordFactor)
	dlat := lat2 - lat1
	dlng := lng2 - lng1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return int32(distance)
}

func toRadians(num float64) float64 {
	return num * math.Pi / float64(180)
}
