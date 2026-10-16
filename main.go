package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		log.Fatal("PORT not found in the environment")
	}

	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		log.Fatal("DATABASE_URL not found in the environment")
	}

	conn, err := sql.Open("postgres", DATABASE_URL)
	if err != nil {
		log.Fatal("Can't connect to the database: ", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(5 * time.Minute)

	// Verify connectivity at startup
	if err = conn.Ping(); err != nil {
		log.Fatal("Database ping failed: ", err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	// Scraper stop channel for graceful shutdown
	stopScraper := make(chan struct{})
	go startScraper(db, 10, time.Minute, stopScraper)

	router := chi.NewRouter()
	v1Router := chi.NewRouter()

	router.Use(loggingMiddleware)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// 404 / 405 handlers
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, 404, "route not found")
	})
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, 405, "method not allowed")
	})

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

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + PORT,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown on SIGINT / SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server is listening on port: %v", PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting the server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutdown signal received, shutting down gracefully...")

	// Stop the scraper goroutine
	close(stopScraper)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited cleanly")
}
