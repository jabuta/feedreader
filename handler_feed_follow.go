package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jabuta/feedreader/internal/database"
)

func (cfg apiConfig) createFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type CreateFeedFollowRequest struct {
		ID uuid.UUID `json:"feed_id"`
	}
	var createFeedFollowRequest CreateFeedFollowRequest
	if err := json.NewDecoder(r.Body).Decode(&createFeedFollowRequest); err != nil {
		respondWithError(w, http.StatusInternalServerError, "can't decode json")
	}
	createdFeedFollow, err := createFeedFollow(cfg.DB, r.Context(), createFeedFollowRequest.ID, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, databaseFeedFollowToMainFeedFollow(createdFeedFollow))
}

func createFeedFollow(db *database.Queries, ctx context.Context, feedId, userId uuid.UUID) (database.FeedFollow, error) {
	return db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feedId,
		UserID:    userId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
}

func (cfg apiConfig) deleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	chi.URLParam(r, "feedFollowID")
	log.Print(chi.URLParam(r, "feedFollowID"))
	feedFollowID, err := uuid.Parse(chi.URLParam(r, "feedFollowID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	deletedFeedFollow, err := cfg.DB.DeleteFeedFollow(r.Context(), feedFollowID)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, databaseFeedFollowToMainFeedFollow(deletedFeedFollow))
}

func (cfg apiConfig) getFeedFollowsUser(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetFeedFollowSUser(r.Context(), user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, http.StatusNotFound, "no feed follows found")
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJson(w, http.StatusOK, databaseFeedFollowsToMainFeedFollows(feedFollows))
}
