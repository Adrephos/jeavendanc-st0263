package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/Adrephos/jeavendanc-st0263/DirectoryServer/proto"
	"github.com/Adrephos/jeavendanc-st0263/DirectoryServer/server"
	"google.golang.org/grpc"
)

var port = flag.Int("port", 50051, "The server port")

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Listening on localhost:%d\n\n", *port)

	grpcServer := grpc.NewServer()
	pb.RegisterDirectoryServiceServer(grpcServer, server.NewServer(time.Minute*5))

	grpcServer.Serve(lis)
}
