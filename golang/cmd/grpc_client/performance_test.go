package main

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	pb "rest_vs_grpc/pkg/polygon"
)

func setupClient() (pb.PolygonServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	client := pb.NewPolygonServiceClient(conn)
	return client, conn, nil
}

func TestPerformance(t *testing.T) {
	clients := 100   // Number of concurrent clients
	requests := 1000 // Number of requests per client

	points := []*pb.Point{
		{X: 0, Y: 0},
		{X: 4, Y: 0},
		{X: 4, Y: 3},
	}

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < clients; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client, conn, err := setupClient()
			if err != nil {
				t.Errorf("could not connect to server: %v", err)
				return
			}
			defer conn.Close()

			for j := 0; j < requests; j++ {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				_, err := client.CalculateArea(ctx, &pb.PolygonRequest{Points: points})
				if err != nil {
					t.Errorf("could not calculate area: %v", err)
					return
				}
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)
	log.Printf("Test completed in %v", duration)
	log.Printf("Total requests: %d", clients*requests)
	log.Printf("Requests per second: %f", float64(clients*requests)/duration.Seconds())
}

func BenchmarkPerformance(b *testing.B) {
	points := []*pb.Point{
		{X: 0, Y: 0},
		{X: 4, Y: 0},
		{X: 4, Y: 3},
	}

	client, conn, err := setupClient()
	if err != nil {
		b.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err := client.CalculateArea(ctx, &pb.PolygonRequest{Points: points})
		if err != nil {
			b.Fatalf("could not calculate area: %v", err)
		}
	}
}
