package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/bid", bidHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
