package api

import (
	"encoding/json"
	"github.com/VenkoChakalov/VoronoiDiagrams/voronoi"
	"net/http"
)

type PointsDto struct {
	X, Y float64
}

func (api *Api) handleVoronoi(w http.ResponseWriter, r *http.Request) error {
	var points []PointsDto
	data := r.PostFormValue("data")
	err := json.Unmarshal([]byte(data), &points)
	vor := voronoi.NewVoronoi(mapPoints(points))
	vor.Process()
	println("ads")
	return err
}

func mapPoints(points []PointsDto) []voronoi.Point {
	voronoiPoints := []voronoi.Point{}
	for _, point := range points {
		voronoiPoints = append(voronoiPoints, voronoi.Point{
			X: float64(point.X),
			Y: float64(point.Y),
		})
	}
	return voronoiPoints
}
