package main

import (
	"log"
	"net/http"
	"one_stop/services"
)

func main() {
	http.HandleFunc("/outlook", services.OutlookHandler);
	http.HandleFunc("/teams", services.TeamsHandler);
	http.HandleFunc("/gitlab", services.GitlabHandler);

	log.Fatal(http.ListenAndServe(":8080", nil))

}

