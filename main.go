package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// Define the address of the TCP server
	serverAddr := "10.54.2.161:53150" // Change this to your server's address
	// Connect to the TCP server
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Error connecting to server: %v\n", err)
	}
	defer conn.Close()

	fmt.Printf("Connected to server at %s\n", serverAddr)

	hexString := "41470102fff20000000000000000000000010000000118010000000000000000"

	// Convert the hex string to a byte slice
	rawRequest, err := hex.DecodeString(hexString)
	if err != nil {
		log.Fatalf("Error decoding hex string: %v\n", err)
	}

	// Send the raw byte request
	_, err = conn.Write(rawRequest)
	if err != nil {
		log.Fatalf("Error sending data: %v\n", err)
	}

	fmt.Println("Raw byte request sent")

	// Optionally, read the response from the server
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatalf("Error reading response: %v\n", err)
	}

	// Print the server's response
	fmt.Printf("Server response: %s\n", string(buffer[:n]))
}
