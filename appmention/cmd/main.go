package main

import (
	"example.com/appmention"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/appmention", appmention.AppMention)
	log.Printf("Listening on localhost:%v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}