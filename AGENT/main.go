package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	pb "agent/proto"
	"google.golang.org/grpc"
)

// Function to compute SHA-256 checksum of a file
func computeChecksum(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

type server struct {
	pb.UnimplementedConfigServiceServer
}

// gRPC function to receive YAML files with checksum validation
func (s *server) SendConfig(stream pb.ConfigService_SendConfigServer) error {
	fmt.Println("Receiving config files...")

	// Ensure the directory exists
	savePath := "./config"
	err := os.MkdirAll(savePath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Response{Status: "Files received successfully"})
		}
		if err != nil {
			return err
		}

		// Compute checksum of received YAML file
		calculatedChecksum := computeChecksum(req.Content)

		// Compare received checksum with calculated one
		if req.Checksum != calculatedChecksum {
			return fmt.Errorf("Checksum mismatch for file %s: expected %s, got %s", req.Filename, req.Checksum, calculatedChecksum)
		}

		// Save the YAML file
		filePath := filepath.Join(savePath, req.Filename)
		err = os.WriteFile(filePath, req.Content, 0644)
		if err != nil {
			return err
		}

		fmt.Printf("✔️ Saved file: %s (Checksum: %s) ✅\n", filePath, calculatedChecksum)
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterConfigServiceServer(s, &server{})

	fmt.Println("gRPC Agent listening on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
