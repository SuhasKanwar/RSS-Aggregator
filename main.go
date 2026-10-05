package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("PORT not found in the enviornment")
	}

	router := chi.NewRouter()
	v1Router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Mount("/v1", v1Router)

	v1Router.Get("/healthz", readinessHandler)

	app := &http.Server{
		Handler: router,
		Addr:    ":" + PORT,
	}

	log.Printf("Server is listening on port: %v", PORT)
	err := app.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
