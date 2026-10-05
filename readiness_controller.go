package main

import "net/http"

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{})
}
