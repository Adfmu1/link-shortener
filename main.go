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
		slog.Error("an error has occured when loading .env", "error", err.Error)
	}

	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	connection, err := sql.Open("postgres", dbUrl)
	if err != nil {
		slog.Error("error while connecting to DB", "error", err.Error())
		return
	}
	defer connection.Close()
	queries := database.New(connection)

	app.DB = queries
	app.loggr = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	mux := http.NewServeMux()

	// GET handlers
	mux.HandleFunc("GET /api/health/", healthHandler)
	mux.HandleFunc("GET /{code}", app.redirectHandler)
	mux.HandleFunc("GET /stats/{code}", app.codeStats)

	// POST handlers
	mux.HandleFunc("POST /api/shorten", app.postCode)

	// fileserver with frontend
	mux.Handle("GET /main/", http.StripPrefix("/main/", http.FileServer(http.Dir("./frontend"))))

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	app.loggr.Info("started server", "port", port)
	err = srv.ListenAndServe()
	if err != nil {
		app.loggr.Error("an error has occured when loading .env", "error", err.Error)
	}
}
