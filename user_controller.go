package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SuhasKanwar/RSS-Aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responsdWithError(w, 400, fmt.Sprint("Error parsing JSON:", err))
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		responsdWithError(w, 400, fmt.Sprint("Couldn't create user:", err))
		return
	}
	respondWithJSON(w, 201, databaseUserToUser(user))
}

func (apiCfg *apiConfig) getUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) getPostsForUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})

	if err != nil {
		responsdWithError(w, 400, fmt.Sprintf("couldn't get posts: %v", err))
	}

	respondWithJSON(w, 200, databasePostsToPosts(posts))
}
