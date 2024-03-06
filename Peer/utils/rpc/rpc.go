package utils_rpc

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Adrephos/jeavendanc-st0263/Peer/client"
	pb "github.com/Adrephos/jeavendanc-st0263/Peer/proto"
	"github.com/Adrephos/jeavendanc-st0263/Peer/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getNetIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
		return "localhost"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func StartServer(port int, dir string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Listening on %s:%d\n\n", getNetIP(), port)

	grpcServer := grpc.NewServer()
	pb.RegisterPeerServer(grpcServer, server.NewServer(dir))

	go grpcServer.Serve(lis)
}

func StartClient(serverAddr string, peerName string, dir string, port int) (PClient *client.PeerClient, conn *grpc.ClientConn) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("fail to dial: %v", err)
	}
	c := pb.NewDirectoryServiceClient(conn)
	PClient = client.NewPClient(peerName, "La vida", dir, fmt.Sprintf("%s:%d", getNetIP(), port), time.Minute*5, c)

	PClient.RegisterClient()

	return PClient, conn
}
