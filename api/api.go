package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Api struct {
	Config ApiConfig
}

func NewApi(config ApiConfig) Api {
	return Api{Config: config}
}

func (api *Api) Handle(router *mux.Router) {
	router.Handle("/voronoi", api.handler(api.handleVoronoi)).Methods("POST")
}

func (api *Api) handler(f func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		beginTime := time.Now()
		r.Body = http.MaxBytesReader(w, r.Body, 100*1024*1024)
		w.Header().Set("Content-Type", "application/json")
		if err := f(w, r); err == nil {
			fmt.Printf("Error at: %s", err)
		}
		defer func() {
			endTime := time.Now()
			fmt.Printf("Time to process the request: %d", endTime.Second()-beginTime.Second())
		}()
	})
}
