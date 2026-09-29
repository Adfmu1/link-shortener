package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	database "github.com/Adfmu1/link_shortener/go_sql"
	"go.rtnl.ai/x/randstr"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from root!")
}

func (app application) postCode(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := reqLink{}
	err := decoder.Decode(&params)
	if err != nil {
		app.loggr.Error("Error while decoding data", slog.String("ERROR", err.Error()))
		respondWithError(w, r, http.StatusBadRequest, "Bad request")
		return
	}

	shortCode, err := app.DB.CheckIfCodeExistsFromUrl(r.Context(), params.URL)
	if err != nil && err != sql.ErrNoRows {
		app.loggr.Error("Error while querying DB", slog.String("ERROR", err.Error()), slog.Int("Status code", http.StatusInternalServerError))
		respondWithError(w, r, http.StatusInternalServerError, "Server error")
		return
	} else if shortCode != "" {
		app.loggr.Info("Shortcode exists for this URL", slog.Int("Status code", http.StatusOK), slog.String("Code", shortCode))
		respondWithJSON(w, r, http.StatusOK, shortenedLink{
			URL: shortCode,
		})
		return
	}

	dbResp, err := app.DB.InsertCode(r.Context(), database.InsertCodeParams{
		Code: randstr.Generate(5, charset),
		Url:  params.URL,
	})

	for err != nil {
		dbResp, err = app.DB.InsertCode(r.Context(), database.InsertCodeParams{
			Code: randstr.Generate(5, charset),
			Url:  params.URL,
		})
	}

	app.loggr.Info("Response logged")
	respondWithJSON(w, r, http.StatusCreated, dbResp)
}
