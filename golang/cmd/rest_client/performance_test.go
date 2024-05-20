package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	rs "rest_vs_grpc/pkg/polygon_rest"
	"sync"
	"testing"
	"time"
)

func TestPerformance(t *testing.T) {
	clients := 100   // Number of concurrent clients
	requests := 1000 // Number of requests per client

	points := []rs.Point{
		{X: 0, Y: 0},
		{X: 4, Y: 0},
		{X: 4, Y: 3},
	}

	var wg sync.WaitGroup
	start := time.Now()

	// Custom HTTP client with keep-alive and timeout settings
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        20000,
		MaxIdleConnsPerHost: 10000,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	var successCount, failureCount int64
	var mu sync.Mutex

	for i := 0; i < clients; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := rs.PolygonRequest{Points: points}
			reqBody, err := json.Marshal(req)
			if err != nil {
				t.Errorf("could not marshal request: %v", err)
				return
			}

			for j := 0; j < requests; j++ {
				resp, err := client.Post("http://localhost:5002/api/polygon/calculate-area", "application/json", bytes.NewBuffer(reqBody))
				if err != nil {
					mu.Lock()
					failureCount++
					mu.Unlock()
					t.Errorf("could not send request: %v", err)
					continue
				}
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					mu.Lock()
					successCount++
					mu.Unlock()
				} else {
					mu.Lock()
					failureCount++
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)
	totalRequests := clients * requests
	successRate := float64(successCount) / float64(totalRequests) * 100
	failureRate := float64(failureCount) / float64(totalRequests) * 100

	log.Printf("Test completed in %v", duration)
	log.Printf("Total requests: %d", totalRequests)
	log.Printf("Successful requests: %d", successCount)
	log.Printf("Failed requests: %d", failureCount)
	log.Printf("Success rate: %.2f%%", successRate)
	log.Printf("Failure rate: %.2f%%", failureRate)
	log.Printf("Requests per second: %f", float64(totalRequests)/duration.Seconds())
}

func BenchmarkPerformance(b *testing.B) {
	points := []rs.Point{
		{X: 0, Y: 0},
		{X: 4, Y: 0},
		{X: 4, Y: 3},
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        20000,
		MaxIdleConnsPerHost: 10000,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	var successCount, failureCount int64
	var mu sync.Mutex

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := rs.PolygonRequest{Points: points}
		reqBody, err := json.Marshal(req)
		if err != nil {
			b.Fatalf("could not marshal request: %v", err)
		}

		resp, err := client.Post("http://localhost:8080/calculate-area", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			mu.Lock()
			failureCount++
			mu.Unlock()
			b.Fatalf("could not send request: %v", err)
		}
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			mu.Lock()
			successCount++
			mu.Unlock()
		} else {
			mu.Lock()
			failureCount++
			mu.Unlock()
		}
	}

	totalRequests := b.N
	successRate := float64(successCount) / float64(totalRequests) * 100
	failureRate := float64(failureCount) / float64(totalRequests) * 100

	log.Printf("Benchmark completed")
	log.Printf("Total requests: %d", totalRequests)
	log.Printf("Successful requests: %d", successCount)
	log.Printf("Failed requests: %d", failureCount)
	log.Printf("Success rate: %.2f%%", successRate)
	log.Printf("Failure rate: %.2f%%", failureRate)
}
