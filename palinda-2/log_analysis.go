package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

type KeyValue struct {
	Key   string
	Value int
}

func Map(logLine string, emit func(string, int)) {
	parts := strings.Fields(logLine)
	if len(parts) >= 4 {
		url := strings.Trim(parts[3], "[]")
		emit(url, 1)
	}
}

func Reduce(key string, values []int, emit func(string, int)) {
	total := 0
	for _, v := range values {
		total += v
	}
	emit(key, total)
}

func main() {
	// 1. Open the log file
	file, err := os.Open("server_logs.txt")
	if err != nil {
		log.Fatalf("Cannot open file, have you generated logs first?: %v", err)
	}
	defer file.Close()

	// Channels for data flow
	linesChan := make(chan string, 1000)          // For the Producer-Consumer map phase
	intermediateChan := make(chan KeyValue, 5000) // For the Shuffle phase

	var mapWG sync.WaitGroup

	fmt.Println("--- Starting MAP Phase (Worker Pool) ---")

	// THE CONSUMERS: Start a fixed number of Map workers (e.g., 8)
	const numWorkers = 8
	for i := 0; i < numWorkers; i++ {
		mapWG.Add(1)
		go func() {
			defer mapWG.Done()

			// The emit function sends data into our intermediate channel
			emit := func(k string, v int) {
				intermediateChan <- KeyValue{Key: k, Value: v}
			}

			// Read lines from the channel until it gets closed
			for line := range linesChan {
				Map(line, emit)
			}
		}()
	}

	// THE PRODUCER: Read the file line-by-line and send to workers
	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			linesChan <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Fel vid läsning av fil: %v", err)
		}
		close(linesChan) // No more lines, signal workers to stop
	}()

	// Wait for all Map workers to finish, then close the intermediate channel
	go func() {
		mapWG.Wait()
		close(intermediateChan)
	}()

	fmt.Println("--- Starting SHUFFLE Phase ---")

	// SHUFFLE: Group all values by key
	groupedData := make(map[string][]int)
	for kv := range intermediateChan {
		groupedData[kv.Key] = append(groupedData[kv.Key], kv.Value)
	}

	fmt.Println("--- Starting REDUCE Phase ---")

	finalResultChan := make(chan KeyValue, len(groupedData))
	var reduceWG sync.WaitGroup

	// PARALLELIZATION: Spin up a goroutine for every UNIQUE KEY
	for key, values := range groupedData {
		reduceWG.Add(1)
		go func(k string, v []int) {
			defer reduceWG.Done()

			emit := func(k string, v int) {
				finalResultChan <- KeyValue{Key: k, Value: v}
			}

			Reduce(k, v, emit)
		}(key, values)
	}

	// Wait for all reducers to finish
	go func() {
		reduceWG.Wait()
		close(finalResultChan)
	}()

	fmt.Println("--- FINAL OUTPUT ---")
	// Print the final aggregated results
	for result := range finalResultChan {
		fmt.Printf("URL: %-15s Total Hits: %d\n", result.Key, result.Value)
	}
}
