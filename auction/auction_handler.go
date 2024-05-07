package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type BiddingServiceResponse struct {
	AdID     string `json:"ad_id"`
	BidPrice int    `json:"bid_price"`
}

func auctionHandler(w http.ResponseWriter, r *http.Request) {

	// Extract AdPlacementId from request
	adPlacementID := r.URL.Query().Get("ad_placement_id")
	log.Println("Request received for AdPlacementId:", adPlacementID)

	biddingServices := []string{"http://bidding:8081/bid", "http://bidding-2:8081/bid", "http://bidding-3:8081/bid"}

	// create channels to receive bid responses and errors
	bidResponses := make(chan *BiddingServiceResponse, len(biddingServices))
	errors := make(chan error, len(biddingServices))

	var wg sync.WaitGroup

	for _, serviceURL := range biddingServices {
		// call the bidding service in a goroutine, add one to the wait group
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			bidResponse, err := callBiddingService(url)
			if err != nil {
				errors <- err
				return
			}
			log.Println("Received bid response from:", url, "ADid", bidResponse.AdID, "Bid Price:", bidResponse.BidPrice)
			bidResponses <- bidResponse
		}(serviceURL)
	}

	log.Println("Waiting for all responses")
	wg.Wait()
	log.Println("All responses received")

	close(bidResponses)
	close(errors)

	// iterate over the bidResponses channel and select the highest bid
	var bestBid *BiddingServiceResponse

	for bid := range bidResponses {
		if bestBid == nil || bid.BidPrice > bestBid.BidPrice {
			bestBid = bid
		}
	}

	if bestBid == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bestBid)

}

func callBiddingService(url string) (*BiddingServiceResponse, error) {

	client := &http.Client{
		Timeout: time.Millisecond * 200,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bidResponse BiddingServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&bidResponse); err != nil {
		return nil, err
	}

	return &bidResponse, nil
}
