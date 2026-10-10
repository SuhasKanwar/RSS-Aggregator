package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SuhasKanwar/RSS-Aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) createFeedHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responsdWithError(w, 400, fmt.Sprint("Error parsing JSON:", err))
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		responsdWithError(w, 400, fmt.Sprint("Couldn't create feed:", err))
		return
	}
	respondWithJSON(w, 201, databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) getFeedHandler(w http.ResponseWriter, r *http.Request, feed database.Feed) {
	respondWithJSON(w, 200, databaseFeedToFeed(feed))
}
