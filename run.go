package main

import (
	"fmt"
	"github.com/VenkoChakalov/VoronoiDiagrams/api"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func run(app api.Api) {
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	router := mux.NewRouter()

	app.Handle(router.PathPrefix("/api").Subrouter())

	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", app.Config.Port),
		Handler:     cors(router),
		ReadTimeout: 2 * time.Minute,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		println("Error")
	}
}

func main() {
	app := api.NewApi(api.NewApiConfig(8222, 0))
	run(app)
}
