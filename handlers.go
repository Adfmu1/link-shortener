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
	type req struct {
		URL string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := req{}
	err := decoder.Decode(&params)
	if err != nil {
		slog.Error("Error while decoding data", "ERROR", err)
	}

	resp := struct {
		URL string `json:"shortUrl"`
	}{
		URL: "https://shorturl.at/iZkiR",
	}
	data, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(data))
}
