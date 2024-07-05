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
	var (
		ctx, cancel = context.WithCancel(context.Background())
		ttl         = 500 * time.Millisecond
		cap         = 50
		removeWG    sync.WaitGroup
		reqWG       sync.WaitGroup
		n           = 100
		rc          = NewRequestCounter(ttl, cap)
		processTime = time.Millisecond * 5
		ts          = httptest.NewServer(NewAPI(processTime, rc))
	)
	defer ts.Close()

	// Remove expired requests
	removeWG.Add(1)
	go func() {
		defer removeWG.Done()
		rc.RemoveExpired(ctx)
	}()

	// Send n requests
	for i := 0; i < n; i++ {
		reqWG.Add(1)
		go func() {
			defer reqWG.Done()
			resp, err := http.Get(ts.URL)
			if err != nil || resp == nil {
				t.Error("failed to send GET request")
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error("failed to read response body")
				return
			}
			if len(body) == 0 {
				t.Error("response body is empty")
				return
			}

			var respBody serverResponse
			err = json.Unmarshal(body, &respBody)
			if err != nil {
				t.Error("failed to unmarshal response body")
				return
			}
			// assert that the number of requests is less than or equal to cap
			if respBody.Count > cap {
				t.Errorf("expected less than %d requests, got %d", i+1, respBody.Count)
			}
		}()
	}

	reqWG.Wait()
	cancel()
	removeWG.Wait()
}
