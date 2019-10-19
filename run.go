package main

import (
	"flag"
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
	fmt.Printf("Server started listening on port %d...", app.Config.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("Error: %s \n", err)
	}
}

func main() {
	port := flag.Int("port", 8222, "the default port that the app is running is 8222")
	proxyCount := flag.Int("proxies", 8222, "the default count of proxies is 0, it's not implemented")
	flag.Parse()
	app := api.NewApi(api.NewApiConfig(*port, *proxyCount))
	run(app)
}
