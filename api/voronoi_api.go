package api

import (
	"encoding/json"
	"net/http"
)

type VoronoiPoint struct {
	X, Y float64
}

func (api *Api) handleVoronoi(w http.ResponseWriter, r *http.Request) error {
	var points []VoronoiPoint
	err := json.NewDecoder(r.Body).Decode(&points)
	return err
}
