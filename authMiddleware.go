package main

import (
	"net/http"

	"github.com/jabuta/feedreader/internal/auth"
	"github.com/jabuta/feedreader/internal/database"
)

type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (cfg apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetHeaderAuth(r)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		user, err := cfg.DB.GetUserByAPI(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		handler(w, r, user)
	}
}
