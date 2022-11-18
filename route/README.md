With gRPC we can define our service once in .proto file and generate clients and servers in any of gRPC's supported languages, which in turn can be run in environments ranging from servers inside a large data center to your own table - all the complexity of communication between different languages and environments is handled for you by gRPC.

route_guide.pb.go: which contains all the protocol buffer code to populate, serialize, and retrieve request and response message types.
route_guide_grpc.pb.go: which contains the following: 
- An interface type (or stub) for clients to call with the methods defined in the RouteGuide server
- An interface type of servers to implement, also with the methods defined in the RouteGuide service