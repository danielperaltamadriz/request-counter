package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func SaveRequests(l []time.Time) error {
	var timeStringList []string = make([]string, len(l))
	for i := 0; i < len(l); i++ {
		timeStringList[i] = l[i].Format(time.RFC3339)
	}

	err := os.WriteFile(fileName, []byte(strings.Join(timeStringList, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}
	log.Println("requests saved")
	return nil
}

func LoadRequests() []time.Time {
	log.Println("loading requests")

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Printf("failed to read file: %v\n", err)
	}
	if len(data) == 0 {
		return nil
	}

	timeStringList := strings.Split(string(data), "\n")
	var l []time.Time = make([]time.Time, len(timeStringList))
	for i := 0; i < len(timeStringList); i++ {
		t, err := time.Parse(time.RFC3339, timeStringList[i])
		if err != nil {
			log.Printf("failed to parse time: %v\n", err)
			continue
		}
		l[i] = t
	}
	return l
}
