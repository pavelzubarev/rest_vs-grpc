package polygon_rest

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type PolygonRequest struct {
	Points []Point `json:"points"`
}

type PolygonResponse struct {
	Area float64 `json:"area"`
}
