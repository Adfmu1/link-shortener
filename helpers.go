package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type errorResponse struct {
	Error string `json:"error"`
}

func respondWithJSON(rw http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		slog.Error("respondWithJSON: marshal failed", "error", err)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusInternalServerError)
		if _, err = rw.Write([]byte(`{"error":"internal server error"}`)); err != nil {
			slog.ErrorContext(r.Context(), "respondWithJSON: write failed", slog.String("error", err.Error()))
		}
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	if _, err := rw.Write(data); err != nil {
		slog.Error("respondWithJSON: write failed", "error", err)
	}
}

func respondWithError(rw http.ResponseWriter, r *http.Request, code int, msg string) {
	respondWithJSON(rw, r, code, errorResponse{Error: msg})
}

func isUrlValid(uri string) bool {
	slog.Info("checking validity")
	if !isUrlReal(uri) {
		return isUrlReal(uri)
	}
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		slog.Error("given url is not valid",
			slog.String("url", uri),
			slog.String("error", err.Error()))
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		slog.Error("given url scheme is not valid",
			slog.String("url", uri))
		return false
	}

	return u.Host != ""
}

func isUrlReal(uri string) bool {
	respClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := respClient.Get(uri)
	if err != nil {
		slog.Error("issue with url has occured",
			slog.String("error", err.Error()))
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		slog.Error("site from url returned bad code",
			slog.Int("status code", resp.StatusCode))
		return false
	}

	slog.Info("url is real", slog.Int("status code", resp.StatusCode))

	return true
}
