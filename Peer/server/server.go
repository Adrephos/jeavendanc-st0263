package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	pb "github.com/Adrephos/jeavendanc-st0263/Peer/proto"
)

type peerServer struct {
	pb.UnimplementedPeerServer

	dir string
}

// Runs in 8080
func (s *peerServer) Download(ctx context.Context, r *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	postBody, _ := json.Marshal(map[string]string{
		"file": r.File,
		"dir":  s.dir,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:8080/download", "application/json", responseBody)
	if err != nil {
		return &pb.DownloadResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	type respStruct struct {
		Success bool   `json:"success"`
		Data    string `json:"data"`
	}
	respStr := respStruct{}
	json.Unmarshal(body, &respStr)
	return &pb.DownloadResponse{
		File:     r.File,
		Metadata: respStr.Data,
	}, nil
}

// Runs in 8081
func (s *peerServer) List(ctx context.Context, r *pb.ListRequest) (*pb.ListResponse, error) {
	postBody, _ := json.Marshal(map[string]string{
		"dir": s.dir,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:8081/list", "application/json", responseBody)
	if err != nil {
		return &pb.ListResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	type respStruct struct {
		Success bool     `json:"success"`
		Data    []string `json:"data"`
	}
	respStr := respStruct{}
	json.Unmarshal(body, &respStr)
	return &pb.ListResponse{Files: respStr.Data}, nil
}

// Runs in 8082
func (s *peerServer) Upload(ctx context.Context, r *pb.UploadRequest) (*pb.UploadResponse, error) {
	postBody, _ := json.Marshal(map[string]string{
		"file": r.File,
		"dir":  s.dir,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:8082/upload", "application/json", responseBody)
	if err != nil {
		return &pb.UploadResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	type respStruct struct {
		Success bool   `json:"success"`
	}
	respStr := respStruct{}
	json.Unmarshal(body, &respStr)
	return &pb.UploadResponse{Success: respStr.Success}, nil
}

func NewServer(dir string) *peerServer {
	s := &peerServer{dir: dir}
	return s
}
