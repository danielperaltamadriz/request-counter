package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestAddToCounter(t *testing.T) {
	var (
		rc = NewRequestCounter(time.Minute)
		n  = 50
		wg sync.WaitGroup
	)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			rc.AddRequest()
		}()
	}
	wg.Wait()

	if rc.CountRequests() != n {
		t.Errorf("expected %d requests, got %d", n, rc.CountRequests())
	}
}

func TestRemoveExpired(t *testing.T) {
	var (
		expireFirst  = time.Now().Add(time.Millisecond * 10)
		expireSecond = time.Now().Add(time.Millisecond * 10)
		rc           = NewRequestCounter(time.Second, expireFirst, expireSecond)
		ctx, cancel  = context.WithCancel(context.Background())
		wg           sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		rc.RemoveExpired(ctx)
	}()
	if rc.CountRequests() != 2 {
		t.Errorf("expected %d requests, got %d", 1, rc.CountRequests())
	}

	// eventually the request will expire and be removed
	<-time.After(time.Until(expireFirst) + time.Millisecond*100)
	if rc.CountRequests() != 0 {
		t.Errorf("expected %d requests, got %d", 0, rc.CountRequests())
	}
	cancel()
	wg.Wait()
}

func TestLoadRequestsImportsOnlyNotExpired(t *testing.T) {
	rc := NewRequestCounter(time.Millisecond,
		time.Time{},
		time.Now().Add(time.Millisecond*-10),
		time.Now().Add(time.Minute),
	)

	if len(rc.GetRequests()) != 1 {
		t.Errorf("expected %d requests, got %d", 1, len(rc.GetRequests()))
	}
}

func TestLoadRequestsSortsRequests(t *testing.T) {
	rc := NewRequestCounter(time.Millisecond,
		time.Now().Add(time.Minute),
		time.Now().Add(time.Second*10),
	)

	requests := rc.GetRequests()
	if len(requests) < 2 {
		t.Errorf("expected at least %d requests, got %d", 2, len(requests))
	}
	if requests[0].After(requests[1]) {
		t.Errorf("requests are not sorted")
	}
}
