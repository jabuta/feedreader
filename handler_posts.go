package main

import (
	"net/http"
	"strconv"

	"github.com/jabuta/feedreader/internal/database"
)

func (cfg apiConfig) getPostsUser(w http.ResponseWriter, r *http.Request, user database.User) {
	limit, numErr := strconv.Atoi(r.URL.Query().Get("limit"))
	if numErr != nil {
		limit = 10
	}

	dbPosts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	posts := databasePostsToPosts(dbPosts)
	respondWithJson(w, http.StatusOK, posts)
}
