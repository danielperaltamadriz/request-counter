package main

import (
	"context"
	"log"
	"sort"
	"sync"
	"time"
)

// RequestCounter is a concurrent safe request counter
type RequestCounter struct {
	requestList []time.Time
	mu          *sync.Mutex
	exp         time.Duration
}

// NewRequestCounter creates a new RequestCounter
// Optionally it can receive a list of requests to load
func NewRequestCounter(exp time.Duration, preLoadedRequests ...time.Time) *RequestCounter {
	return &RequestCounter{
		requestList: loadRequests(preLoadedRequests),
		exp:         exp,
		mu:          &sync.Mutex{},
	}
}

// GetRequests returns all requests in the request list.
func (rc *RequestCounter) GetRequests() []time.Time {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	return rc.requestList
}

// AddRequest adds a request to the request list.
func (rc *RequestCounter) AddRequest() {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.requestList = append(rc.requestList, time.Now().Add(rc.exp))
}

// CountRequests returns the number of requests in the request list.
func (rc *RequestCounter) CountRequests() int {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	return len(rc.requestList)
}

// RemoveExpired removes expired requests from the request list. It runs until the context is done.
func (rc *RequestCounter) RemoveExpired(ctx context.Context) {
	ticker := time.NewTicker(rc.getNextExpiration())

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			log.Println("context is done")
			return
		case <-ticker.C:
			rc.mu.Lock()
			if len(rc.requestList) > 0 && time.Now().After(rc.requestList[0]) {
				rc.requestList = rc.requestList[1:]
			}
			rc.mu.Unlock()
			ticker.Reset(rc.getNextExpiration())
		}
	}
}

func (rc *RequestCounter) getNextExpiration() time.Duration {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if len(rc.requestList) == 0 {
		return rc.exp
	}
	return time.Until(rc.requestList[0])
}

// LoadRequests loads requests from a slice of time.Time. It only loads requests that are not expired.
func loadRequests(requests []time.Time) []time.Time {
	sort.Slice(requests, func(i, j int) bool {
		return requests[i].Before(requests[j])
	})

	var validRequests []time.Time
	for i := 0; i < len(requests); i++ {
		if time.Now().Before(requests[i]) {
			validRequests = append(validRequests, requests[i])
		}
	}
	return validRequests
}
