package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from root!")
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := reqLink{}
	err := decoder.Decode(&params)
	if err != nil {
		slog.Error("Error while decoding data", "ERROR", err)
		respondWithError(w, http.StatusBadRequest, "Bad request")
		return
	}

	resp := shortenedLink{
		URL: "https://shorturl.at/iZkiR",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Error while decoding data", "ERROR", err)
		respondWithError(w, http.StatusInternalServerError, "Bad request")
		return
	}

	respondWithJSON(w, http.StatusCreated, data)
}
