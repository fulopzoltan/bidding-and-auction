package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/google/uuid"
)

type AdObject struct {
	AdID     uuid.UUID `json:"ad_id"`
	BidPrice int       `json:"bid_price"`
}

func bidHandler(w http.ResponseWriter, r *http.Request) {

	// simulate not bidding and returning no content
	if rand.Intn(10) < 2 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Simulate a random bid price
	bidPrice := rand.Intn(100)
	adObject := AdObject{
		AdID:     uuid.New(), // random adId
		BidPrice: bidPrice,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(adObject)
}

func main() {

	http.HandleFunc("/bid", bidHandler)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
