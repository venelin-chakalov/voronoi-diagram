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
	data := r.PostFormValue("points")
	err := json.Unmarshal([]byte(data), &points)
	if err != nil {
		return err
	}
	vor := voronoi.NewVoronoi(mapPoints(points))
	vor.Process()
	err = json.NewEncoder(w).Encode(vor.Result)
	return err
}

func mapPoints(points []PointsDto) []voronoi.Point {
	var voronoiPoints []voronoi.Point
	for _, point := range points {
		voronoiPoints = append(voronoiPoints, voronoi.Point{
			X: 1.0 * (point.X),
			Y: 1.0 * (point.Y),
		})
	}
	return voronoiPoints
}
