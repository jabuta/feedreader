package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/jabuta/feedreader/internal/database"
)

func isValidURL(toTest string) bool {
	parsedURL, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false // Not a valid URL
	}
	return (parsedURL.Scheme == "http" || parsedURL.Scheme == "https") && parsedURL.Host != ""
}

func (cfg apiConfig) handlerFeedsPost(w http.ResponseWriter, r *http.Request, user database.User) {
	type PostFeed struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	var postFeed PostFeed

	if err := json.NewDecoder(r.Body).Decode(&postFeed); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if postFeed.Name == "" || !isValidURL(postFeed.Url) {
		respondWithError(w, http.StatusBadRequest, "invalid name or url")
		return
	}

	createdFeed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      postFeed.Name,
		Url:       postFeed.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, databaseFeedToMainFeed(createdFeed))
}

func (cfg apiConfig) getFeeds(w http.ResponseWriter, r *http.Request) {
	dbFeeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJson(w, http.StatusOK, databaseFeedsToMainFeeds(dbFeeds))
}
