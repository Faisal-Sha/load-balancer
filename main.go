package main

import (
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	backendServers = []string{
		"http://localhost:8080",
		"http://localhost:8081",
		"http://localhost:8082",
	}
	healthyServers []string
	mutex          sync.RWMutex
	counter        uint64
)

func main() {
	listenPort := ":3000"
	healthCheckInterval := 10 * time.Second

	// Initialize healthy servers
	healthyServers = backendServers
	go startHealthChecks(healthCheckInterval)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request from %s\n", r.RemoteAddr)

		server := getHealthyServer()
		if server == "" {
			http.Error(w, "No healthy servers available", http.StatusServiceUnavailable)
			return
		}

		log.Printf("Forwarding request to %s\n", server)

		resp, err := http.Get(server)
		if err != nil {
			http.Error(w, "Backend server error", http.StatusInternalServerError)
			log.Printf("Error forwarding request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)

		log.Printf("Response from %s: %s\n", server, resp.Status)
	})

	log.Printf("Load Balancer running on port %s\n", listenPort)
	if err := http.ListenAndServe(listenPort, nil); err != nil {
		log.Fatalf("Failed to start load balancer: %v", err)
	}
}

func startHealthChecks(interval time.Duration) {
	for {
		time.Sleep(interval)
		checkHealth()
	}
}

func checkHealth() {
	log.Println("Performing health checks...")
	var newHealthyServers []string
	for _, server := range backendServers {
		resp, err := http.Get(server)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Printf("Server %s is unhealthy\n", server)
			continue
		}
		newHealthyServers = append(newHealthyServers, server)
	}

	mutex.Lock()
	healthyServers = newHealthyServers
	mutex.Unlock()
	log.Printf("Healthy servers: %v\n", healthyServers)
}

func getHealthyServer() string {
	mutex.RLock()
	defer mutex.RUnlock()

	if len(healthyServers) == 0 {
		return ""
	}
	next := atomic.AddUint64(&counter, 1)
	return healthyServers[next%uint64(len(healthyServers))]
}
