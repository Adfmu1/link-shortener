package main

import (
	"fmt"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from root!")
}
