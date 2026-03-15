package collector

import (
	"fmt"      // Used to print logs/events to console
	"io"       // Used to detect EOF when connection closes
	"net"      // Provides networking primitives (TCP listener, connections)
	"strconv"  // Used to convert string ↔ integer
	// "strings"  // (Imported but not used in this file)
	"time"     // Used to measure connection duration
	"sentinelx/models"
)

// StartTCPMonitor starts a TCP server that monitors incoming connections
func StartTCPMonitor(port string) {

	// Start listening on the given TCP port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err) // Stop program if port cannot be opened
	}

	fmt.Println("TCP monitor listening on port", port)

	for {

		// Accept new incoming TCP connection
		conn, err := listener.Accept()
		if err != nil {
			continue // If accept fails, skip and wait for next connection
		}

		// Handle each connection in a separate goroutine
		// This allows multiple connections simultaneously
		go handleConnection(conn)

	}
}

// handleConnection processes one TCP connection
func handleConnection(conn net.Conn) {

		defer conn.Close()

	// Record when connection started
	start := time.Now()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// Get remote client address (example: 192.168.1.5:52344)
	remoteAddr := conn.RemoteAddr().String()

	// Split into host IP and port
	host, portStr, _ := net.SplitHostPort(remoteAddr)

	// Convert port string → integer
	portInt := 0
	if p, err := strconv.Atoi(portStr); err == nil {
		portInt = p
	}

	// -------- CONNECTION OPEN EVENT --------

	// Create security event when connection opens
	openEvent := models.NewSecurityEvent("connection_open")

	openEvent.SourceIP = host      // Client IP
	openEvent.SourcePort = portInt // Client port
	openEvent.Protocol = "TCP"     // Protocol used

	// Convert event to JSON and print
jsonData, err := openEvent.ToJSON()
if err == nil {
    fmt.Println(string(jsonData))
}

	// Buffer to read incoming TCP data
	buffer := make([]byte, 4096)

	// Track how many bytes were transferred
	totalBytes := 0

	for {

    n, err := conn.Read(buffer)

    if n > 0 {
        totalBytes += n
    }

    if err != nil {

        if err != io.EOF {
            // unexpected error
            fmt.Println("connection error:", err)
        }

        break
    }
}

	// Calculate how long connection lasted
	duration := time.Since(start)

	// -------- CONNECTION CLOSE EVENT --------

	// Create event when connection closes
	closeEvent := NewSecurityEvent("connection_close")

	closeEvent.SourceIP = host
	closeEvent.SourcePort = portInt
	closeEvent.Protocol = "TCP"

	// Store extra metadata
	closeEvent.Metadata["duration"] = duration.String()
	closeEvent.Metadata["bytes_transferred"] = strconv.Itoa(totalBytes)

	// Convert event to JSON and print
	jsonData, err = closeEvent.ToJSON()
if err == nil {
	fmt.Println(string(jsonData))
}

	
}