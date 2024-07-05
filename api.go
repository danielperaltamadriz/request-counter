package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type API struct {
	rc             *RequestCounter
	reqProcessTime time.Duration
}

func NewAPI(processTime time.Duration, rc *RequestCounter) *API {
	return &API{
		reqProcessTime: processTime,
		rc:             rc,
	}
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.rc.AddRequest()
	time.Sleep(api.reqProcessTime)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(fmt.Sprintf(`{"request_count": %s}`, strconv.Itoa(api.rc.CountRequests()))))
	if err != nil {
		log.Printf("failed to write response: %v\n", err)
	}
}
