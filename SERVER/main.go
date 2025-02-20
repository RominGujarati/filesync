package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"github.com/fsnotify/fsnotify"
	pb "server/proto"
	"server/db"
	"server/ping_listener"
	"server/routes"
	"server/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func onConfigChange(event fsnotify.Event) {
	eventType := ""

	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
		eventType = "create"
	case event.Op&fsnotify.Write == fsnotify.Write:
		eventType = "modify"
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		eventType = "delete"
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		eventType = "rename"
	}

	if eventType == "" {
		return
	}
	
	fmt.Printf("📂 Event: %s | File: %s\n", eventType, event.Name)
	triggerMainFunction(eventType, event.Name)
}

// Placeholder for the function to trigger
func triggerMainFunction(eventType string, filename string) {
	fmt.Println("Triggered main function due to config change.")
	// send all the ./config dir with all the files to agent via a grpc call
	SendConfigToAgent(eventType, filename)
}

func SendConfigToAgent(eventType string, filename string) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewConfigServiceClient(conn)
	stream, err := client.SendConfig(context.Background())
	if err != nil {
		log.Printf("Error creating stream: %v", err)
		return
	}

	configPath := filename
	content :=  []byte("")

	if eventType != "delete" {

		content, err = ioutil.ReadFile(configPath)
		if err != nil {
			fmt.Println("Error reading config file:", err)
		}
		
	}

	// Compute checksum for YAML file
	checksum := utils.ComputeChecksum(content)

	req := &pb.ConfigFile{
		Filename: filename,
		Content:  content,
		Checksum: checksum,
		Eventtype: eventType,
	}

	if err := stream.Send(req); err != nil {
		fmt.Println("Error sending file:", err)
	}
	fmt.Printf("📤 Sent file: %s (Checksum: %s)\n", configPath, checksum)

	if err != nil {
		fmt.Println("Error reading config directory:", err)
	}

	stream.CloseAndRecv()
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

	err = watcher.Add(configPath + "/jobs")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening for changes in the /config and /config/jobs folders...")

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

func main() {
	db.ConnectDB()
	r := gin.Default()
	// Register authentication routes
	routes.AuthRoutes(r)
	routes.FileRoutes(r)
	routes.ProtectedRoutes(r)
	go r.Run(":8080")
	go ping_listener.StartPingListener()
	listenForChanges()
}
