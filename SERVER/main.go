package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	// "net"
	"os"
	"path/filepath"
	"github.com/fsnotify/fsnotify"
	"crypto/sha256"
	"encoding/hex"
	
	pb "server/proto"
	// "server/models"
	"server/db"
	// "gopkg.in/yaml.v2"
	"google.golang.org/grpc"
	"server/ping_listener"
)

// Function to trigger on file changes
func onConfigChange(event fsnotify.Event) {
	fmt.Printf("Config changed: %s, Event: %s\n", event.Name, event.Op)
	// Call your main function or any specific function
	triggerMainFunction()
}

// Placeholder for the function to trigger
func triggerMainFunction() {
	fmt.Println("Triggered main function due to config change.")
	// send all the ./config dir with all the files to agent via a grpc call
	SendConfigToAgent()
}


// Function to compute SHA-256 checksum
func computeChecksum(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func SendConfigToAgent() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewConfigServiceClient(conn)
	stream, err := client.SendConfig(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	configPath := "./config"
	err = filepath.Walk(configPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".yaml" { // Send only YAML files
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// Compute checksum for YAML file
			checksum := computeChecksum(content)

			req := &pb.ConfigFile{
				Filename: filepath.Base(path),
				Content:  content,
				Checksum: checksum,
			}

			if err := stream.Send(req); err != nil {
				return err
			}
			fmt.Printf("📤 Sent file: %s (Checksum: %s)\n", path, checksum)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error reading config directory: %v", err)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println("✔️ Agent response:", resp.Status)
}
func listenForChanges() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	configPath := "./config"

	err = watcher.Add(configPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening for changes in the /config folder...")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			onConfigChange(event) // Trigger function when an event occurs

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}

// Config struct for reading from YAML
// type Config struct {
// 	Server struct {
// 		IP   string `yaml:"ip"`
// 		Port string `yaml:"port"`
// 	} `yaml:"server"`
// }

// LoadConfig reads the YAML file
// func loadConfig(path string) (*Config, error) {
// 	file, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var config Config
// 	err = yaml.Unmarshal(file, &config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &config, nil
// }


func main() {
	db.ConnectDB()
	go ping_listener.StartPingListener()
	listenForChanges()
}
