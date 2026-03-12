package main

import (
	"fmt"                 // Used to print messages to console
	"net/http"            // Go standard library for building HTTP servers
	"sentinelx/collector" // Import our collector package which contains HTTPCollector middleware
)

func main() {

	// Create a new HTTP request multiplexer (router)
	// It is used to register URL routes like /login, /home, etc.
	mux := http.NewServeMux()

	// Register a handler for the "/login" endpoint
	// When someone visits http://localhost:8080/login this function runs
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		// Send a simple response back to the client
		fmt.Fprintf(w, "Login page")
	})

	// Wrap the router with the HTTPCollector middleware
	// This means every request first passes through HTTPCollector
	// where a SecurityEvent will be created
	handler := collector.HTTPCollector(mux)

	// Print message to console so we know server started
	fmt.Println("Server running on :8080")

	// Start the HTTP server on port 8080
	// All requests go through the middleware handler
	http.ListenAndServe(":8080", handler)
}