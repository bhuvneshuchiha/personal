package services

import (
	"fmt"
	"net/http"
)

func TeamsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Fprintf(w, "Hello")
	} else {
		fmt.Fprintf(w, "Not a post request")
	}
}
