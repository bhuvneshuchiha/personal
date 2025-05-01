package services

import (
	"fmt"
	"net/http"
)

func GitlabHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Fprintf(w, "Hello")
	} else {
		fmt.Fprintf(w, "Not a post request")
	}
}
