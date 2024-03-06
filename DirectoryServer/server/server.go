package server

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	pb "github.com/Adrephos/jeavendanc-st0263/DirectoryServer/proto"
)

const (
	UP   = "UP"
	DOWN = "DOWN"
)

type node struct {
	url           string
	state         string
	lastKeepalive time.Time
	files         map[string]bool
}

type directoryService struct {
	pb.UnimplementedDirectoryServiceServer

	mu               sync.Mutex
	db               map[string]*node
	keepaliveTimeout time.Duration
}

func (s *directoryService) addFiles(nodeName string, files ...string) error {
	node, ok := s.db[nodeName]
	s.mu.Lock()
	if !ok {
		return errors.New("node not found")
	}
	for _, file := range files {
		node.files[file] = true
	}
	s.mu.Unlock()
	return nil
}

func (s *directoryService) updateNode(name string, url string, files ...string) error {
	newNode := &node{url: url, state: UP, files: make(map[string]bool)}
	s.mu.Lock()
	if _, ok := s.db[name]; !ok {
		s.db[name] = newNode
	} else {
		s.db[name].url = url
	}
	s.mu.Unlock()
	s.addFiles(name, files...)
	return nil
}

func (s *directoryService) refreshKeepalive(name string) (time.Time, error) {
	if _, ok := s.db[name]; !ok {
		return time.Now(), errors.New("node not found")
	}
	s.mu.Lock()
	s.db[name].lastKeepalive = time.Now()
	s.mu.Unlock()
	return s.db[name].lastKeepalive, nil
}

func (s *directoryService) changeNodeState(name string, state string) error {
	if _, ok := s.db[name]; !ok {
		return errors.New("node not found")
	}
	s.mu.Lock()
	s.db[name].state = state
	s.mu.Unlock()
	return nil
}

func (s *directoryService) deleteFiles(name string) error {
	s.mu.Lock()
	if _, ok := s.db[name]; !ok {
		return errors.New("node not found")
	}
	s.db[name].files = make(map[string]bool)
	s.mu.Unlock()
	return nil
}

func (s *directoryService) searchFile(file string) ([]*pb.SearchResponse_FileInfo, error) {
	var nodes []*pb.SearchResponse_FileInfo
	err := errors.New("file not found")
	for name, node := range s.db {
		valid, _ := s.validKeepalive(name)
		if node.state == DOWN || !valid {
			s.changeNodeState(name, DOWN)
			continue
		}
		if _, ok := node.files[file]; ok {
			nodes = append(nodes, &pb.SearchResponse_FileInfo{
				Node: name,
				Url:  node.url,
			})
			err = nil
		}
	}
	return nodes, err
}

func (s *directoryService) validKeepalive(name string) (bool, error) {
	if _, ok := s.db[name]; !ok {
		return false, errors.New("node not found")
	}
	elapsed := time.Since(s.db[name].lastKeepalive)
	if elapsed > s.keepaliveTimeout {
		return false, nil
	}
	return true, nil
}

func (s *directoryService) Login(ctx context.Context, b *pb.LoginRequest) (*pb.LoginResponse, error) {
	if b.Name == "" || b.Password == "" {
		return &pb.LoginResponse{}, errors.New("no credentials")
	}
	s.updateNode(b.Name, "")
	s.changeNodeState(b.Name, UP)
	s.refreshKeepalive(b.Name)

	log.Printf("New login with name %s\n", b.Name)
	return &pb.LoginResponse{
		Token:   os.Getenv("TOKEN"),
		Success: true,
	}, nil
}

func (s *directoryService) Logout(ctx context.Context, b *pb.NodeName) (*pb.LogoutResponse, error) {
	if b.Name == "" {
		return &pb.LogoutResponse{}, errors.New("no credentials")
	}
	err := s.deleteFiles(b.Name)
	err = s.changeNodeState(b.Name, DOWN)
	if err != nil {
		return &pb.LogoutResponse{}, errors.New("node not available")
	}

	log.Printf("new logout with name %s\n", b.Name)
	return &pb.LogoutResponse{
		Message: "logout successful",
		Success: true,
	}, nil
}

func (s *directoryService) Search(ctx context.Context, b *pb.SearchRequest) (*pb.SearchResponse, error) {
	nodes, err := s.searchFile(b.File)
	if err != nil {
		return &pb.SearchResponse{}, err
	}

	log.Printf("searching for file %s\n", b.File)
	return &pb.SearchResponse{
		Response: nodes,
		Success:  true,
	}, nil
}

func (s *directoryService) Keepalive(ctx context.Context, b *pb.NodeName) (*pb.KeepaliveResponse, error) {
	if b.Name == "" {
		return &pb.KeepaliveResponse{}, errors.New("no credentials")
	}
	time, err := s.refreshKeepalive(b.Name)
	err = s.changeNodeState(b.Name, UP)
	if err != nil {
		return &pb.KeepaliveResponse{}, err
	}

	log.Printf("node %s sent keepalive\n", b.Name)
	return &pb.KeepaliveResponse{
		Response: &pb.KeepaliveResponse_NodeInfo{
			Name:          b.Name,
			LastKeepalive: time.String(),
		},
		Success: true,
	}, nil
}

func (s *directoryService) Index(ctx context.Context, b *pb.IndexRequest) (*pb.IndexResponse, error) {
	if b.Name == "" {
		return &pb.IndexResponse{}, errors.New("no credentials")
	}
	s.updateNode(b.Name, b.Url, b.Files...)
	s.refreshKeepalive(b.Name)

	log.Printf("node %s indexed new files %v\n", b.Name, b.Files)
	return &pb.IndexResponse{Success: true}, nil
}

func (s *directoryService) GetPeers(ctx context.Context, b *pb.PeersRequest) (*pb.PeersResponse, error) {
	var response []*pb.PeersResponse_NodeInfo
	for name, node := range s.db {
		nodeInfo := &pb.PeersResponse_NodeInfo{
			Name: name, Url: node.url}
		response = append(response, nodeInfo)
	}
	return &pb.PeersResponse{
		Response: response,
		Success: true,
	}, nil
}

func NewServer(keepaliveTimeout time.Duration) *directoryService {
	s := &directoryService{
		db:               make(map[string]*node),
		keepaliveTimeout: keepaliveTimeout}
	return s
}
