package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	// Start the load balancer on port 80
	listenPort := ":3000"
	backendServer := "http://localhost:8080" // Single backend server

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request from %s\n", r.RemoteAddr)

		// Forward the request to the backend server
		resp, err := http.Get(backendServer)
		if err != nil {
			http.Error(w, "Backend server error", http.StatusInternalServerError)
			log.Printf("Error forwarding request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// Copy backend response to the client
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)

		log.Printf("Response from backend server: %s\n", resp.Status)
	})

	log.Printf("Load Balancer running on port %s\n", listenPort)
	if err := http.ListenAndServe(listenPort, nil); err != nil {
		log.Fatalf("Failed to start load balancer: %v", err)
	}
}
