package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SuhasKanwar/RSS-Aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("PORT not found in the enviornment")
	}

	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		log.Fatal("DATABASE_URL not found in the enviornment")
	}

	conn, conn_err := sql.Open("postgres", DATABASE_URL)
	if conn_err != nil {
		log.Fatal("Can't connect to the database: ", conn_err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	go startScraper(db, 10, time.Minute)

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

	v1Router.Post("/users", apiCfg.createUserHandler)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.getUserHandler))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.createFeedHandler))
	v1Router.Get("/feeds", apiCfg.getFeedsHandler)

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.getPostsForUserHandler))

	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.createFeedFollowHandler))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.getFeedFollowsHandler))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.deleteFeedFollowHandler))

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
