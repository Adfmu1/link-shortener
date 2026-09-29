package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	database "github.com/Adfmu1/link_shortener/go_sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type application struct {
	DB    *database.Queries
	loggr *slog.Logger
}

var app application

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("An error has occured when loading .env", "error", err.Error)
	}

	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	connection, err := sql.Open("postgres", dbUrl)
	if err != nil {
		slog.Error("Error while connecting to DB", "Error", err.Error())
		return
	}
	defer connection.Close()
	queries := database.New(connection)

	app.DB = queries
	app.loggr = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/", rootHandler)
	mux.HandleFunc("POST /api/shorten", app.postCode)

	mux.Handle("GET /", http.FileServer(http.Dir("./frontend")))

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	app.loggr.Info("Started server", "Port", port)
	err = srv.ListenAndServe()
	if err != nil {
		app.loggr.Error("An error has occured when loading .env", "error", err.Error)
	}
}
