package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	database "github.com/Adfmu1/link_shortener/go_sql"
	"go.rtnl.ai/x/randstr"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, r, http.StatusOK, struct {
		Response string `json:"response"`
	}{
		Response: "healthy",
	})
}

func (app application) codeStats(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	data, err := app.DB.GetDataFromCode(r.Context(), code)
	if errors.Is(err, sql.ErrNoRows) {
		app.loggr.Error("code doesnt exist in DB",
			slog.String("code", code))
		respondWithError(w, r, http.StatusNotFound, "code doesnt exist")
		return
	}
	if err != nil {
		app.loggr.Error("an error has occured when accessing the code data",
			slog.String("code", code))
		respondWithError(w, r, http.StatusNotFound, "code doesnt exist")
		return
	}

	app.loggr.InfoContext(r.Context(), "accesed short url data",
		slog.String("Code", code))

	respondWithJSON(w, r, http.StatusOK, data)
}

func (app application) postCode(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := reqLink{}
	err := decoder.Decode(&params)
	if err != nil {
		app.loggr.Error("error while decoding data",
			slog.String("error", err.Error()))
		respondWithError(w, r, http.StatusBadRequest, "bad request")
		return
	}

	if !isUrlValid(params.URL) {
		respondWithError(w, r, http.StatusBadRequest, "URI not valid")
		return
	}

	shortCode, err := app.DB.CheckIfCodeExistsFromUrl(r.Context(), params.URL)
	if err != nil && err != sql.ErrNoRows {
		app.loggr.Error("error while querying DB",
			slog.String("error", err.Error()),
			slog.Int("status code", http.StatusInternalServerError))
		respondWithError(w, r, http.StatusInternalServerError, "Server error")
		return
	} else if shortCode != "" {
		app.loggr.Info("shortcode exists for this URL",
			slog.Int("status code", http.StatusOK),
			slog.String("code", shortCode))
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

	app.loggr.Info("response logged")
	respondWithJSON(w, r, http.StatusCreated, dbResp)
}

func (app application) redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	data, err := app.DB.GetDataFromCode(r.Context(), code)
	if err != nil {
		respondWithError(w, r, http.StatusNotFound, "code doesnt exist")
		return
	}
	app.loggr.InfoContext(r.Context(), "accesed short url",
		slog.String("Code", code),
		slog.String("Url", data.Url))

	err = app.DB.IncrementClicksFromCode(r.Context(), code)
	if err != nil {
		respondWithError(w, r, http.StatusNotFound, "code doesnt exist")
		app.loggr.Error("error with incrementing the click count for given code",
			slog.String("code", code))
		return
	}
	http.Redirect(w, r, data.Url, http.StatusFound)
}
