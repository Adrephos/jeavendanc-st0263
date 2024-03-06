package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"slices"
	"strings"
	"time"

	pb "github.com/Adrephos/jeavendanc-st0263/Peer/proto"
	utils_files "github.com/Adrephos/jeavendanc-st0263/Peer/utils/files"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var lastSentFiles []string

type PeerClient struct {
	peerName         string
	token            string
	directory        string
	url              string
	keepaliveTimeout time.Duration
	dsClient         pb.DirectoryServiceClient
}

func (c *PeerClient) sendKeepalive() {
	go func() {
		for {
			c.Index()
		}
	}()
	for {
		c.dsClient.Keepalive(context.Background(),
			&pb.NodeName{Name: c.peerName},
		)
		time.Sleep(c.keepaliveTimeout)
	}
}

func (c *PeerClient) peerConection(url string) (pb.PeerClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("could not connect to peer")
		return nil, nil, errors.New("could not connect to peer")
	}
	pc := pb.NewPeerClient(conn)

	return pc, conn, nil
}

func (c *PeerClient) Index() {
	files := utils_files.DirFiles(c.directory)
	if slices.Equal(files, lastSentFiles) {
		return
	}
	lastSentFiles = files
	c.dsClient.Index(context.Background(),
		&pb.IndexRequest{
			Name:  c.peerName,
			Url:   c.url,
			Files: files,
		},
	)
}

func (c *PeerClient) RegisterClient() error {
	_, err := c.dsClient.Login(context.Background(),
		&pb.LoginRequest{Name: c.peerName, Password: c.token},
	)
	if err != nil {
		return err
	}
	go c.sendKeepalive()

	log.Println("client registered to Directory Server with name", c.peerName)
	return nil
}

func (c *PeerClient) Search(file string) (string, error) {
	r, err := c.dsClient.Search(context.Background(),
		&pb.SearchRequest{File: strings.TrimSpace(file)})
	if err != nil {
		return "", err
	}
	nodes, _ := json.MarshalIndent(r.Response, "", " ")
	log.Printf("file found on peers:\n%s\n", string(nodes))
	i := rand.IntN(len(r.Response))
	log.Printf("using %s to get the file", r.Response[i].Node)
	return r.Response[i].Url, nil
}

func (c *PeerClient) Logout() error {
	r, err := c.dsClient.Logout(context.Background(), &pb.NodeName{Name: c.peerName})
	log.Println(r.Message)
	return err
}

func (c *PeerClient) GetPeers() (string, error) {
	r, err := c.dsClient.GetPeers(context.Background(), &pb.PeersRequest{})
	if err != nil {
		return "", err
	}
	nodes, _ := json.MarshalIndent(r.Response, "", " ")
	log.Printf("online peers:\n%s\n", string(nodes))

	return string(nodes), nil
}

func (c *PeerClient) Download(file string, url string) error {
	client, conn, err := c.peerConection(url)
	defer conn.Close()
	if err != nil {
		return err
	}

	r, err := client.Download(context.Background(),
		&pb.DownloadRequest{File: file})
	if err != nil {
		return err
	}

	log.Printf("new file dowloaded, metadata:\n %s\n", r.Metadata)
	utils_files.CreateFile(c.directory, r.File)
	log.Printf("created %s in %s\n", file, c.directory)
	return nil
}

func (c *PeerClient) List(url string) error {
	client, conn, err := c.peerConection(url)
	defer conn.Close()
	if err != nil {
		return err
	}

	r, err := client.List(context.Background(),
		&pb.ListRequest{})
	if err != nil {
		return err
	}

	files, _ := json.MarshalIndent(r.Files, "", " ")
	log.Printf("got this files from: %s \n%s\n", url, string(files))
	return nil
}

func (c *PeerClient) Upload(file string, url string) error {
	client, conn, err := c.peerConection(url)
	defer conn.Close()
	if err != nil {
		return err
	}

	r, err := client.Upload(context.Background(),
		&pb.UploadRequest{File: file})
	if err != nil {
		return err
	}

	logText := fmt.Sprintf("uploaded file %s to %s", file, url)
	if !r.Success {
		logText = fmt.Sprintf("failed to upload file %s to %s", file, url)
	}
	log.Println(logText)
	return nil
}

func NewPClient(
	name string,
	token string,
	directory string,
	url string,
	keepaliveTimeout time.Duration,
	dsclient pb.DirectoryServiceClient) *PeerClient {

	lastSentFiles = append(lastSentFiles, " ")
	return &PeerClient{
		name,
		token,
		directory,
		url,
		keepaliveTimeout,
		dsclient,
	}
}
