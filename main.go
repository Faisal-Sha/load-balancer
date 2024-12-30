package main

import (
	"io"
	"log"
	"net/http"
	"sync/atomic"
)

var (
	backendServers = []string{
		"http://localhost:8080",
		"http://localhost:8081",
	}
	counter uint64
)

func main() {
	listenPort := ":80"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request from %s\n", r.RemoteAddr)

		// Get the next server using Round Robin
		server := getNextServer()
		log.Printf("Forwarding request to %s\n", server)

		// Forward request
		resp, err := http.Get(server)
		if err != nil {
			http.Error(w, "Backend server error", http.StatusInternalServerError)
			log.Printf("Error forwarding request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// Copy backend response to the client
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)

		log.Printf("Response from %s: %s\n", server, resp.Status)
	})

	log.Printf("Load Balancer running on port %s\n", listenPort)
	if err := http.ListenAndServe(listenPort, nil); err != nil {
		log.Fatalf("Failed to start load balancer: %v", err)
	}
}

func getNextServer() string {
	// Round Robin algorithm
	next := atomic.AddUint64(&counter, 1)
	return backendServers[next%uint64(len(backendServers))]
}
