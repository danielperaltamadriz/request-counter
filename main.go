package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	fileName = "requests.csv"

	defaultTTLSec = 60
	defaultPort   = "8080"
)

func main() {
	// Parse expiration time and port from environment variables
	var err error
	ttlSec, err := strconv.Atoi(os.Getenv("TTL_SEC"))
	if err != nil {
		log.Printf("failed to parse TTL_SEC as number: %s\n", err)
	}
	ttl := time.Duration(ttlSec) * time.Second
	if ttlSec <= 0 {
		log.Printf("invalid TTL_SEC: %d, using default value: %d\n", ttlSec, defaultTTLSec)
		ttl = time.Duration(defaultTTLSec) * time.Second
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("PORT is not set, using default value: %s\n", defaultPort)
		port = defaultPort
	}

	// Create concurrent safe request counter and load requests from file
	rc := NewRequestCounter(ttl, LoadRequests()...)

	// Create context and wait group
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Start goroutine to remove expired requests
	wg.Add(1)
	go func() {
		defer wg.Done()
		rc.RemoveExpired(ctx)
	}()

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: NewAPI(rc),
	}

	// Start HTTP server
	wg.Add(1)
	done := make(chan struct{})
	go func() {
		defer wg.Done()
		log.Printf("Listening on :%s\n", port)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("failed to listen and serve: %v\n", err)
		}
		close(done)
	}()

	// Create signal channel if SIGINT or SIGTERM is received
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal or server to be done
	select {
	// If signal is received, shutdown server
	case <-quit:
		log.Println("Shutdown Server")
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("failed to shutdown server: %v\n", err)
		}

	case <-done:
		log.Println("Server is done")
	}

	// Cancel context and wait for goroutines to finish
	cancel()
	wg.Wait()

	// Save requests to file
	err = SaveRequests(rc.GetRequests())
	if err != nil {
		log.Printf("failed to save requests: %v\n", err)
	}
}
