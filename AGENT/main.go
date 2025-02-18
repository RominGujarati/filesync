package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	// "time"
	// "context"
	"gopkg.in/yaml.v2"

	"agent/utils"

	pb "agent/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Send ping notification to the server
func sendPingNotification(message string) {

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
	conn, err := net.Dial("udp", "localhost:6000")
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

// PingServer sends the agent's IP and timestamp to the server every 1 hour in a separate goroutine
// func pingServer(client pb.ConfigServiceClient, agentIP string) {
// 	for {
// 		timestamp := time.Now().Format(time.RFC3339)

// 		_, err := client.Ping(context.Background(), &pb.PingRequest{
// 			AgentIp:   agentIP,
// 			Timestamp: timestamp,
// 		})
// 		if err != nil {
// 			log.Printf("Failed to send ping: %v", err)
// 		} else {
// 			log.Printf("Ping sent to server: %s at %s", agentIP, timestamp)
// 		}

// 		time.Sleep(1 * time.Hour) // Wait 1 hour before the next ping
// 	}
// }


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
		calculatedChecksum := utils.ComputeChecksum(req.Content)

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

		sendPingNotification("Config file received")

		fmt.Printf("✔️ Saved file: %s (Checksum: %s) ✅\n", filePath, calculatedChecksum)
	}

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


		// Create a gRPC client to send pings
		// conn, err := grpc.Dial(serverIP+":"+serverPort, grpc.WithInsecure())
		// if err != nil {
		// 	log.Fatalf("Failed to connect to server: %v", err)
		// }
		// defer conn.Close()
		// client := pb.NewConfigServiceClient(conn)
	
		// // Get agent's IP address
		// agentIP, err := getLocalIP()
		// if err != nil {
		// 	log.Fatalf("Failed to get agent IP: %v", err)
		// }
	
		// Start pinging in a separate goroutine
		// go sendPingNotification("Agent started")

}
