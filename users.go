package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jabuta/feedreader/internal/database"
)

func (cfg apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type userCreate struct {
		Name string `json:"name"`
	}
	var reqBody userCreate
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	createUserArgs := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      reqBody.Name,
	}
	user, err := cfg.DB.CreateUser(r.Context(), createUserArgs)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, user)
}
