package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("An error has occured when loading .env", "error", err.Error)
	}

	port := os.Getenv("PORT")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/", rootHandler)
	mux.HandleFunc("POST /api/shorten", postHandler)

	mux.Handle("GET /", http.FileServer(http.Dir("./frontend")))

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("Started server", "Port", port)
	err = srv.ListenAndServe()
	if err != nil {
		slog.Error("An error has occured when loading .env", "error", err.Error)
	}
}
