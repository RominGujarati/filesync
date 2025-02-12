package ping_listener

import (
	"fmt"
	"log"
	"net"

	"server/db"
	"server/models"
)

func StartPingListener() {
	addr := ":6000"
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("Error resolving UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("Error starting UDP server: %v", err)
	}
	defer conn.Close()

	log.Printf("📡 Ping listener started on %s", addr)

	buffer := make([]byte, 1024)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("📩 Received ping: %s from %s\n", message, remoteAddr.String())

		// Save ping to DB
		ping := models.Ping{
			Status:  "Received",
			Details: message,
			IP:      remoteAddr.String(),
		}

		result := db.DB.Create(&ping)
		if result.Error != nil {
			log.Printf("❌ Failed to save ping: %v", result.Error)
		} else {
			log.Println("✅ Ping saved to database")
		}
	}
}
