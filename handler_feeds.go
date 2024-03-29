package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/jabuta/feedreader/internal/database"
)

func (cfg apiConfig) handlerFeedsPost(w http.ResponseWriter, r *http.Request, user database.User) {
	type PostFeed struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	type createFeedResponse struct {
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
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
		ID:     uuid.New(),
		Name:   postFeed.Name,
		Url:    postFeed.Url,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:     uuid.New(),
		FeedID: createdFeed.ID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, createFeedResponse{
		Feed:       databaseFeedToFeed(createdFeed),
		FeedFollow: databaseFeedFollowToFeedFollow(feedFollow),
	})
}

func (cfg apiConfig) handlerFeedsGet(w http.ResponseWriter, r *http.Request) {
	dbFeeds, err := cfg.DB.GetFeeds(r.Context())
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, http.StatusNotFound, "no feeds found")
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, databaseFeedsToFeeds(dbFeeds))
}

func isValidURL(toTest string) bool {
	parsedURL, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false // Not a valid URL
	}
	return (parsedURL.Scheme == "http" || parsedURL.Scheme == "https") && parsedURL.Host != ""
}
