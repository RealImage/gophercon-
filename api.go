package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/time/rate"
)

const API_ENDPOINT = "http://data.moviebuff.com/"

// NewClient with a ratelimiter
func NewClient(rl *rate.Limiter) *HTTPClient {
	c := &HTTPClient{
		client:      http.DefaultClient,
		RateLimiter: rl,
	}
	return c
}

// A wrapper over client.Do() method for Rate limiting.
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	err := c.RateLimiter.Wait(req.Context())
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Generic Function to Fetch Person|Movie Details
func FetchEntityDetails[T Entity](url string) (*T, error) {
	req, err := http.NewRequest(http.MethodGet, API_ENDPOINT+url, nil)
	if err != nil {
		return nil, err
	}

	// Reduce the following limit in case of http.StatusTooManyRequests
	rl := rate.NewLimiter(rate.Every(1*time.Second), 10000) // 10000 requests per second
	client := NewClient(rl)

	res, err := client.Do(req)

	switch true {
	case err != nil:
		log.Println("Error occurred")
		return nil, err

	case res.StatusCode != http.StatusOK:
		return nil, fmt.Errorf("%d: error occurred", res.StatusCode)

	// In case of DoS prevention from the CDN, reduce the rate limit and try again
	case res.StatusCode == http.StatusTooManyRequests:
		log.Println("Reduce Rate Limit and Try Again!")
		os.Exit(1)
	}

	var entity T
	if err := json.NewDecoder(res.Body).Decode(&entity); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return &entity, nil
}

func GetNames(parentURL string, personURL string, movieURL string) (string, string, string) {
	parent, err := FetchEntityDetails[Person](parentURL)
	if err != nil {
		log.Println(err)
	}

	person, err := FetchEntityDetails[Person](personURL)
	if err != nil {
		log.Println(err)
	}

	movie, err := FetchEntityDetails[Movie](movieURL)
	if err != nil {
		log.Println(err)
	}
	return parent.Name, person.Name, movie.Name
}
