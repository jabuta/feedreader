package main

import "net/http"

func v1OkHandler(w http.ResponseWriter, r *http.Request) {
	type OkResponse struct {
		Status string `json:"status"`
	}
	respondWithJson(w, http.StatusOK, OkResponse{
		Status: "ok",
	})
}

func v1NotOkHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
