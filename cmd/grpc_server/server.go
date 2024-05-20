package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math"
	"net"
	pb "rest_vs_grpc/pkg/polygon"
)

type server struct {
	pb.UnimplementedPolygonServiceServer
}

func (s *server) CalculateArea(ctx context.Context, req *pb.PolygonRequest) (*pb.PolygonResponse, error) {
	points := req.GetPoints()
	if len(points) < 3 {
		return nil, fmt.Errorf("a polygon must have at least 3 points")
	}

	area := 0.0
	n := len(points)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += points[i].GetX() * points[j].GetY()
		area -= points[j].GetX() * points[i].GetY()
	}
	area = math.Abs(area) / 2.0

	return &pb.PolygonResponse{Area: area}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPolygonServiceServer(s, &server{})

	log.Printf("grpc_server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
