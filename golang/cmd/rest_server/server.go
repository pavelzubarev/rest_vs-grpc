package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	rs "rest_vs_grpc/pkg/polygon_rest"
	"time"
)

func calculateArea(points []rs.Point) float64 {
	if len(points) < 3 {
		return 0.0
	}

	area := 0.0
	n := len(points)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += points[i].X * points[j].Y
		area -= points[j].X * points[i].Y
	}
	area = math.Abs(area) / 2.0

	return area
}

func calculateAreaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req rs.PolygonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	area := calculateArea(req.Points)
	resp := rs.PolygonResponse{Area: area}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/api/polygon/calculate-area", calculateAreaHandler)

	server := &http.Server{
		Addr:           ":5002",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	log.Fatal(server.ListenAndServe())
}
