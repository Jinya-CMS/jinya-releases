package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})

	log.Println("Serving at localhost:8090...")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
