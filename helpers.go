package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func respondWithJSON(rw http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		slog.Error("respondWithJSON: marshal failed", "error", err)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{"error":"internal server error"}`))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	if _, err := rw.Write(data); err != nil {
		slog.Error("respondWithJSON: write failed", "error", err)
	}
}

func respondWithError(rw http.ResponseWriter, code int, msg string) {
	if code >= http.StatusInternalServerError {
		slog.Error("Error has occured:", "Status code", code, "Error", msg)
	}
	respondWithJSON(rw, code, errorResponse{Error: msg})
}
