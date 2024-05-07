package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Server up and running on port 8080")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Auction Service, please use /auction endpoint to get the best bid."))
	})
	http.HandleFunc("/auction", auctionHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
