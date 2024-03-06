package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	utils_cmd "github.com/Adrephos/jeavendanc-st0263/Peer/utils/cmd"
	"github.com/Adrephos/jeavendanc-st0263/Peer/utils/microservices"
	utils_rpc "github.com/Adrephos/jeavendanc-st0263/Peer/utils/rpc"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	peerName   = flag.String("name", "peer-1", "The name to identify the peer")
	peerPort   = flag.Int("port", 50052, "The port the peer will be running on")
	directory  = flag.String("dir", "/files/", "The port the peer will be running on")
	download   = flag.String("download", "./microservices/download/download", "Path to download microservice executable")
	upload     = flag.String("upload", "./microservices/upload/upload", "Path to upload microservice executable")
	list       = flag.String("list", "./microservices/list/list", "Path to list microservice executable")
)

func main() {
	flag.Parse()
	// Set name and port with env variable
	env := os.Getenv("PEER_NAME")
	if env != "" {
		flag.Set("name", env)
	}
	env = os.Getenv("DIR_SERVER_ADDR")
	if env != "" {
		flag.Set("addr", env)
	}
	env = os.Getenv("PORT")
	if env != "" {
		flag.Set("port", env)
	}
	// Make dir an absoulute path
	abs, err := filepath.Abs(*directory)
	if err != nil {
		log.Fatalln("not a valid directory")
	}
	directory = &abs

	cmd := exec.Command("mkdir", *directory)
	cmd.Run()

	// Server
	utils_rpc.StartServer(*peerPort, *directory)

	// Start microservices
	microservices.StartMicroservices(*download, *list, *upload)

	// Client
	PClient, conn := utils_rpc.StartClient(*serverAddr, *peerName, *directory, *peerPort)
	utils_cmd.CommandLine(PClient)
	PClient.Logout()
	conn.Close()
}
