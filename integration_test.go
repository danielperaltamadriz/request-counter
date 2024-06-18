package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type serverResponse struct {
	Count int `json:"request_count"`
}

func TestIntegration(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rc := NewRequestCounter(10 * time.Millisecond)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		rc.RemoveExpired(ctx)
	}()

	ts := httptest.NewServer(NewAPI(rc))
	defer ts.Close()

	requestTime := time.Now().Add(time.Millisecond * 20)
	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Error("failed to send GET request")
	}
	if resp != nil {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error("failed to read response body")
		}
		if len(body) > 0 {
			var respBody serverResponse
			err = json.Unmarshal(body, &respBody)
			if err != nil {
				t.Error("failed to unmarshal response body")
			}
			if respBody.Count != 1 {
				t.Errorf("expected 1 request, got %d", respBody.Count)
			}
		}
	}
	<-time.After(time.Until(requestTime))
	resp, err = http.Get(ts.URL)
	if err != nil {
		t.Error("failed to send GET request")
	}
	if resp != nil {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error("failed to read response body")
		}
		if len(body) > 0 {
			var respBody serverResponse
			err = json.Unmarshal(body, &respBody)
			if err != nil {
				t.Error("failed to unmarshal response body")
			}
			if respBody.Count != 1 {
				t.Errorf("expected 1 requests, got %d", respBody.Count)
			}
		}
	}
	cancel()
	wg.Wait()
}
