package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	rs "rest_vs_grpc/pkg/polygon_rest"
)

func main() {
	points := []rs.Point{
		{X: 0, Y: 0},
		{X: 4, Y: 0},
		{X: 4, Y: 3},
	}

	req := rs.PolygonRequest{Points: points}
	reqBody, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("could not marshal request: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/calculate-area", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	var res rs.PolygonResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Fatalf("could not decode response: %v", err)
	}

	log.Printf("Area: %f", res.Area)
}
