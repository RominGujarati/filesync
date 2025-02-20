package main

import (
	pb "agent/proto"
	"agent/utils"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// Send ping notification to the server
func sendPingNotification(message string, serverIP string, serverPort string) {

	// Load configuration
	// config, err := loadConfig("C:/Users/bhard/OneDrive/Documents/sites/kumaran/filesync/AGENT/config/server-config.yaml")

	// if err != nil {
	// 	log.Fatalf("Failed to load config: %v", err)
	// }

	// serverIP := config.Server.IP
	// serverPort := config.Server.Port

	// log.Printf("Starting gRPC Agent on %s:%s...", serverIP, serverPort)

	// conn, err := net.Dial("udp", serverIP+":"+serverPort)
	// localhost :6000
	// conn, err := net.Dial("udp", "localhost:6000")
	conn, err := net.Dial("udp", serverIP+":"+serverPort)
	if err != nil {
		fmt.Printf("Failed to send ping: %v\n", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Error sending ping: %v\n", err)
	} else {
		fmt.Println("📡 Ping sent successfully")
	}
}

type server struct {
	pb.UnimplementedConfigServiceServer
}

// Config struct for reading from YAML
type Config struct {
	Server struct {
		IP   string `yaml:"ip"`
		Port string `yaml:"port"`
	} `yaml:"server"`
}

// LoadConfig reads the YAML file
func loadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetLocalIP retrieves the agent's local IP address
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}
	return "", fmt.Errorf("could not determine local IP address")
}

// gRPC function to receive YAML files with checksum validation and handle deletions
func (s *server) SendConfig(stream pb.ConfigService_SendConfigServer) error {
	fmt.Println("Receiving config files...")
	req, err := stream.Recv()

	fileName := req.Filename
	tmpdir := strings.Split(fileName, "\\")
	fileName = tmpdir[len(tmpdir)-1]
	tmpdir = tmpdir[:len(tmpdir)-1]
	dir := strings.Join(tmpdir, "\\")

	// Ensure the directory exists
	savePath := "./"
	filePath := filepath.Join(savePath, dir, fileName)
	direrr := os.MkdirAll(dir, 0755)
	if direrr != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	if err == io.EOF {
		return stream.SendAndClose(&pb.Response{Status: "Files received successfully"})
	}
	if err != nil {
		return err
	}

	// Handle file deletion event
	if req.Eventtype == "delete" {
		err = os.Remove(filePath)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete file %s: %v", filePath, err)
		}
		fmt.Printf("🗑️ Deleted file: %s\n", filePath)
		return nil
	}

	// Compute checksum of received YAML file
	calculatedChecksum := utils.ComputeChecksum(req.Content)

	// Compare received checksum with calculated one
	if req.Checksum != calculatedChecksum {
		return fmt.Errorf("Checksum mismatch for file %s: expected %s, got %s", req.Filename, req.Checksum, calculatedChecksum)
	}

	// Save the YAML file
	err = os.WriteFile(filePath, req.Content, 0644)
	if err != nil {
		return err
	}

	// TO DO FIX:
	sendPingNotification("Config file received", "localhost", "6000")

	fmt.Printf("✔️ Saved file: %s (Checksum: %s) ✅\n", filePath, calculatedChecksum)

	return nil

}

func main() {
	// Load configuration
	config, err := loadConfig("C:/Users/bhard/OneDrive/Documents/sites/kumaran/filesync/AGENT/config/server-config.yaml")

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	serverIP := config.Server.IP
	serverPort := config.Server.Port

	log.Printf("Starting gRPC Agent on %s:%s...", serverIP, serverPort)

	// Start gRPC Server
	lis, err := net.Listen("tcp", serverIP+":"+serverPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterConfigServiceServer(s, &server{})
	reflection.Register(s)

	fmt.Println("gRPC Agent listening on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
