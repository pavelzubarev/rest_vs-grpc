package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "rest_vs_grpc/pkg/polygon"
	rs "rest_vs_grpc/pkg/polygon_rest"
)

func setupGrpcClient() (pb.PolygonServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	client := pb.NewPolygonServiceClient(conn)
	return client, conn, nil
}

func setupHttpClient() *http.Client {
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
	return client
}

func benchmarkGrpcPerformance(points []*pb.Point, clients, requests int) {
	client, conn, err := setupGrpcClient()
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	var successCount, failureCount int64
	var mu sync.Mutex

	start := time.Now()

	for i := 0; i < clients; i++ {
		wg.Add(1)
		go func(clientIndex int) {
			defer wg.Done()
			for j := 0; j < requests; j++ {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				_, err := client.CalculateArea(ctx, &pb.PolygonRequest{Points: points})
				mu.Lock()
				if err != nil {
					failureCount++
				} else {
					successCount++
				}
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)
	totalRequests := clients * requests
	successRate := float64(successCount) / float64(totalRequests) * 100
	failureRate := float64(failureCount) / float64(totalRequests) * 100

	log.Printf("gRPC Test completed in %v", duration)
	log.Printf("gRPC Total requests: %d", totalRequests)
	log.Printf("gRPC Successful requests: %d", successCount)
	log.Printf("gRPC Failed requests: %d", failureCount)
	log.Printf("gRPC Success rate: %.2f%%", successRate)
	log.Printf("gRPC Failure rate: %.2f%%", failureRate)
	log.Printf("gRPC Requests per second: %f", float64(totalRequests)/duration.Seconds())
}

func benchmarkRestPerformance(points []rs.Point, clients, requests int) {
	client := setupHttpClient()

	var wg sync.WaitGroup
	var successCount, failureCount int64
	var mu sync.Mutex

	start := time.Now()

	for i := 0; i < clients; i++ {
		wg.Add(1)
		go func(clientIndex int) {
			defer wg.Done()
			req := rs.PolygonRequest{Points: points}
			reqBody, err := json.Marshal(req)
			if err != nil {
				log.Printf("could not marshal request: %v", err)
				return
			}

			for j := 0; j < requests; j++ {
				resp, err := client.Post("http://localhost:8080/calculate-area", "application/json", bytes.NewBuffer(reqBody))
				mu.Lock()
				if err != nil || resp.StatusCode != http.StatusOK {
					failureCount++
				} else {
					successCount++
				}
				mu.Unlock()
				if resp != nil {
					resp.Body.Close()
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)
	totalRequests := clients * requests
	successRate := float64(successCount) / float64(totalRequests) * 100
	failureRate := float64(failureCount) / float64(totalRequests) * 100

	log.Printf("REST Test completed in %v", duration)
	log.Printf("REST Total requests: %d", totalRequests)
	log.Printf("REST Successful requests: %d", successCount)
	log.Printf("REST Failed requests: %d", failureCount)
	log.Printf("REST Success rate: %.2f%%", successRate)
	log.Printf("REST Failure rate: %.2f%%", failureRate)
	log.Printf("REST Requests per second: %f", float64(totalRequests)/duration.Seconds())
}

func TestComparePerformance(t *testing.T) {
	grpcPoints := []*pb.Point{
		{X: 0, Y: 0},
		{X: 4, Y: 0},
		{X: 4, Y: 3},
	}
	restPoints := []rs.Point{
		{X: 0, Y: 0},
		{X: 4, Y: 0},
		{X: 4, Y: 3},
	}

	clients := 100   // Number of concurrent clients
	requests := 1000 // Number of requests per client

	t.Run("gRPC Performance", func(t *testing.T) {
		benchmarkGrpcPerformance(grpcPoints, clients, requests)
	})

	t.Run("REST Performance", func(t *testing.T) {
		benchmarkRestPerformance(restPoints, clients, requests)
	})
}
