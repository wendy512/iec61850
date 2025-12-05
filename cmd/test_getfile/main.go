package main

import (
	"fmt"
	"log"
	"os"

	"github.com/boeboe/iec61850"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: test_getfile <host:port> <filename>")
		fmt.Println("Example: test_getfile 192.168.100.57:102 index.html")
		os.Exit(1)
	}

	hostPort := os.Args[1]
	filename := os.Args[2]

	// Parse host and port
	host := "localhost"
	port := 102
	fmt.Sscanf(hostPort, "%s:%d", &host, &port)

	settings := iec61850.Settings{
		Host:           host,
		Port:           port,
		ConnectTimeout: 10000, // 10 seconds
		RequestTimeout: 60000, // 60 seconds - increased for large files
	}

	client, err := iec61850.NewClient(settings)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	fmt.Printf("Connected to %s:%d\n", host, port)
	fmt.Printf("Downloading file: %s\n", filename)

	// Test GetFile with the fixed implementation
	data, err := client.GetFile(filename)
	if err != nil {
		log.Fatalf("Failed to get file: %v", err)
	}

	fmt.Printf("Successfully downloaded %d bytes\n", len(data))

	// Save to local file
	outputFile := "downloaded_" + filename
	err = os.WriteFile(outputFile, data, 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Printf("Saved to: %s\n", outputFile)
}
