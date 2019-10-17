package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func (api *Api) run() {
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	router := mux.NewRouter()
	api.Init(router.PathPrefix("/api").Subrouter())

	s := &http.Server{
		Addr:        fmt.Sprintf(":%d", api.Config.Port),
		Handler:     cors(router),
		ReadTimeout: 2 * time.Minute,
	}

}
