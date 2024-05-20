package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "rest_vs_grpc/pkg/polygon"
	"time"
)

func main() {
	conn, err := grpc.NewClient("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPolygonServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	points := []*pb.Point{
		{X: 0, Y: 0},
		{X: 4, Y: 0},
		{X: 4, Y: 3},
	}

	r, err := c.CalculateArea(ctx, &pb.PolygonRequest{Points: points})
	if err != nil {
		log.Fatalf("could not calculate area: %v", err)
	}
	log.Printf("Area: %f", r.GetArea())
}
